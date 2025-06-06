package centraluniversity.app.booking.services;

import java.time.LocalDate;
import java.time.LocalTime;
import java.time.ZoneId;
import java.time.format.DateTimeFormatter;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;

import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;

import centraluniversity.app.booking.models.booking.*;
import centraluniversity.app.booking.models.exception.HttpStatusException;
import centraluniversity.app.booking.models.rooms.Room;
import centraluniversity.app.booking.models.user.GetUserDto;
import centraluniversity.app.booking.pb.Role;
import centraluniversity.app.booking.repositories.BookingRepository;
import lombok.RequiredArgsConstructor;

@Service
@RequiredArgsConstructor
public class BookingService {

    private final AuthService authService;
    private final RoomService roomService;
    private final BookingRepository bookingRepository;
    private final AdminUserService adminUserService;

    /**
     * Проверка, либо что запрос отправлен администратором, либо что пользователь редактирует именно свою запись.
     * @param userId
     * @param email
     * @param password
     * @param isAdmin
     * @throws HttpStatusException NOT_FOUND (почта не найдена), UNAUTHORIZED (пароль не совпадает), FORBIDDEN (доступ запрещен)
     */
    private void auth(Integer userId, String email, String password, boolean isAdmin) throws HttpStatusException {
        GetUserDto user = authService.getUserByEmail(email, password);
        if (isAdmin) {
            if (user.getRole() == Role.USER) {
                throw new HttpStatusException(HttpStatus.FORBIDDEN, "Редактировать чужие брони может только администратор");
            }
            adminUserService.getUserById(userId, email, password);
        } else {
            if (user.getId() != userId) {
                throw new HttpStatusException(HttpStatus.FORBIDDEN, "Редактировать чужие брони может только администратор");
            }
        }
    }

    /**
     * Получение информации о брони по его id.
     * @param id
     * @return полная информация о брони
     * @throws HttpStatusException NOT_FOUND (бронь не найдена)
     */
    private Booking getBookingById(Integer id) throws HttpStatusException {
        Optional<Booking> booking = bookingRepository.findById(id);
        if (booking.isEmpty()) {
            throw new HttpStatusException(HttpStatus.NOT_FOUND, String.format("Бронь с id = %d не существует", id));
        }
        return booking.get();
    }

    /**
     * Преобразование List < Booking > в List < BookingDto >.
     * @param bookings List < Booking >
     * @return List < BookingDto >
     */
    public List<BookingDto> parseToBookingDtoList(List<Booking> bookings) {
        List<BookingDto> result = new ArrayList<>();

        for (int i = 0; i < bookings.size(); i++) {
            Booking bookingNum = bookings.get(i);
            Room room = roomService.getRoomById(bookingNum.getRoomId());
            BookingDto bookingDto = new BookingDto(
                bookingNum.getUserId(), 
                bookingNum.getRoomId(), 
                bookingNum.getBookingDate(), 
                bookingNum.getBookingStart(), 
                bookingNum.getBookingEnd(), 
                room
            );
            result.add(bookingDto);
        }
        return result;
    }

    /**
     * Валидация даты и времени бронирования.
     * @param bookingDate
     * @param bookingStart
     * @param bookingEnd
     * @throws HttpStatusException BAD_REQUEST 
     */
    public static void validateDateTime(LocalDate bookingDate, LocalTime bookingStart, LocalTime bookingEnd) throws HttpStatusException {
        if (bookingDate.isBefore(LocalDate.now(ZoneId.of("Europe/Moscow")))) {
            throw new HttpStatusException(HttpStatus.BAD_REQUEST, "Нельзя забронировать аудиторию ранее сегодняшнего дня");
        }

        if (bookingStart.isAfter(bookingEnd) || bookingStart.equals(bookingEnd)) {
            throw new HttpStatusException(HttpStatus.BAD_REQUEST, "Время начала должно быть строго раньше времени окончания");
        }
    }
    
    /**
     * Получение всех записей по аудитории за определенную дату
     * @param roomId
     * @param bookingDate
     * @return
     */
    public DateBookingDto getTimesByDate(Integer roomId, LocalDate bookingDate) {
        List<Booking> bookings = bookingRepository.findByRoomIdAndBookingDateOrderByBookingStart(roomId, bookingDate);

        List<TimeBookingDto> bookingTimes = new ArrayList<>();
        for (int i = 0; i < bookings.size(); i++) {
            TimeBookingDto time = new TimeBookingDto();
            time.setBookingStart(bookings.get(i).getBookingStart());
            time.setBookingEnd(bookings.get(i).getBookingEnd());
            bookingTimes.add(time);
        }
        return new DateBookingDto(bookingDate, bookingTimes);
    }

    // TODO: unit-tests
    /**
     * Проверка на пересечения времени бронирования.
     * @param times - список всех интервалов за день
     * @param bookingStart
     * @param bookingEnd
     */
    public static void validateTimesByOneDay(List<TimeBookingDto> times, LocalTime bookingStart, LocalTime bookingEnd) {
        DateTimeFormatter timeFormatter = DateTimeFormatter.ofPattern("HH:mm:ss");

        for (int i = 0; i < times.size(); i++) {
            TimeBookingDto time = times.get(i);
            boolean startIsBad = bookingStart.isAfter(time.getBookingStart()) || bookingStart.equals(time.getBookingStart());
            boolean endIsBad = bookingEnd.isBefore(time.getBookingEnd()) || bookingEnd.equals(time.getBookingEnd());
            if (startIsBad && endIsBad) {
                throw new HttpStatusException(HttpStatus.BAD_REQUEST, String.format("Уже существует бронь с %s по %s", 
                    time.getBookingStart().format(timeFormatter), time.getBookingEnd().format(timeFormatter)));
            }
        }
    }

    /**
     * Сохранение новой брони пользователем или администратором.
     * @param booking - информация о брони
     * @param isAdmin
     * @throws HttpStatusException BAD_REQUEST (время брони занято), NOT_FOUND (почта/ауд. не найдена), UNAUTHORIZED (пароль не совпадает), FORBIDDEN (доступ запрещен)
     */
    public void createBooking(CreateBookingDto booking, boolean isAdmin) throws HttpStatusException {

        System.out.println(booking.getEmail());
        System.out.println(booking.getPassword());
        System.out.println(booking.getRoomId());
        System.out.println(booking.getUserId());
        System.out.println(booking.getBookingDate());
        System.out.println(booking.getBookingEnd());
        System.out.println(booking.getBookingStart());
        
        auth(booking.getUserId(), booking.getEmail(), booking.getPassword(), isAdmin);
        roomService.getRoomById(booking.getRoomId());

        validateDateTime(booking.getBookingDate(), booking.getBookingStart(), booking.getBookingEnd());
        List<TimeBookingDto> bookingTimes = getTimesByDate(booking.getRoomId(), booking.getBookingDate()).getBookingTimes();
        validateTimesByOneDay(bookingTimes, booking.getBookingStart(), booking.getBookingEnd());

        Booking bookingSql = new Booking();
        bookingSql.setUserId(booking.getUserId());
        bookingSql.setRoomId(booking.getRoomId());
        bookingSql.setBookingDate(booking.getBookingDate());
        bookingSql.setBookingStart(booking.getBookingStart());
        bookingSql.setBookingEnd(booking.getBookingEnd());

        bookingRepository.save(bookingSql);
    }

    /**
     * Получение всех-всех-всех броней администратором.
     * @param adminEmail
     * @param adminPassword
     * @return пустой список или список с полной информацией о всех бронях
     * @throws HttpStatusException NOT_FOUND (почта не найдена), UNAUTHORIZED (пароль не совпадает), FORBIDDEN (доступ запрещен)
     */
    public List<BookingDto> getAllBookings(String adminEmail, String adminPassword) throws HttpStatusException {
        GetUserDto user = authService.getUserByEmail(adminEmail, adminPassword);
        if (user.getRole() == Role.USER) {
            throw new HttpStatusException(HttpStatus.FORBIDDEN, "Доступ запрещен");
        }
        List<Booking> bookings = bookingRepository.findAll();
        return parseToBookingDtoList(bookings);
    }

    /**
     * Получение пользователем всех собственных броней
     * @param email
     * @param password
     * @return пустой список или список с полной информацией о всех бронях
     * @throws HttpStatusException NOT_FOUND (почта не найдена), UNAUTHORIZED (пароль не совпадает)
     */
    public List<BookingDto> getAllBookingsByEmail(String email, String password) throws HttpStatusException {
        GetUserDto user = authService.getUserByEmail(email, password);
        List<Booking> bookings = bookingRepository.findByUserId(user.getId());
        return parseToBookingDtoList(bookings);
    }

    /**
     * Обновление информации о бронировании пользователем или администратором.
     * @param id
     * @param booking - полная информация о брони
     * @param isAdmin
     * @throws HttpStatusException BAD_REQUEST (время брони занято), NOT_FOUND (почта/ауд. не найдена), UNAUTHORIZED (пароль не совпадает), FORBIDDEN (доступ запрещен)
     */
    public void updateBooking(Integer id, UpdateBookingDto booking, boolean isAdmin) throws HttpStatusException {
        
        auth(booking.getUserId(), booking.getEmail(), booking.getPassword(), isAdmin);
        roomService.getRoomById(booking.getRoomId());
        Booking bookingSql = getBookingById(id);

        LocalTime oldBookingStart = bookingSql.getBookingStart();

        if (booking.getBookingDate() != null) {
            bookingSql.setBookingDate(booking.getBookingDate());
        }
        if (booking.getBookingStart() != null) {
            bookingSql.setBookingStart(booking.getBookingStart());
        }
        if (booking.getBookingEnd() != null) {
            bookingSql.setBookingEnd(booking.getBookingEnd());
        }

        validateDateTime(bookingSql.getBookingDate(), bookingSql.getBookingStart(), bookingSql.getBookingEnd());
        List<TimeBookingDto> bookingTimes = getTimesByDate(booking.getRoomId(), booking.getBookingDate()).getBookingTimes();
        for (int i = 0; i < bookingTimes.size(); i++) {
            if (bookingTimes.get(i).getBookingStart().equals(oldBookingStart)) {
                bookingTimes.remove(i);
                break;
            }
        }
        validateTimesByOneDay(bookingTimes, oldBookingStart, oldBookingStart);
        

        bookingRepository.save(bookingSql);
    } 

    /**
     * Отмена бронирования пользователем или администратором.
     * @param id
     * @param email
     * @param password
     * @param isAdmin
     * @throws HttpStatusException NOT_FOUND (почта/бронь не найдена), UNAUTHORIZED (пароль не совпадает), FORBIDDEN (доступ запрещен)
     */
    public void deleteBooking(Integer id, String email, String password, boolean isAdmin) throws HttpStatusException {
        auth(id, email, password, isAdmin);
        bookingRepository.deleteById(id);
    }
}
