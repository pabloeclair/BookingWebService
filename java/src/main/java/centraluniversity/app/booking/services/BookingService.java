package centraluniversity.app.booking.services;

import java.util.ArrayList;
import java.util.List;
import java.util.Optional;

import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;

import centraluniversity.app.booking.models.booking.Booking;
import centraluniversity.app.booking.models.booking.BookingDto;
import centraluniversity.app.booking.models.booking.CreateBookingDto;
import centraluniversity.app.booking.models.booking.UpdateBookingDto;
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

    // TODO: unit-test
    /**
     * Преобразование List < Booking > в List < BookingDto >.
     * @param bookings List < Booking >
     * @return List < BookingDto >
     */
    private List<BookingDto> parseToBookingDtoList(List<Booking> bookings) {
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
    
    // TODO: сделать обработку дат
    /**
     * Сохранение новой брони пользователем или администратором.
     * @param booking - информация о брони
     * @param isAdmin
     * @throws HttpStatusException //TODO: обновить ошибки
     */
    public void createBooking(CreateBookingDto booking, boolean isAdmin) throws HttpStatusException {
        
        auth(booking.getUserId(), booking.getEmail(), booking.getPassword(), isAdmin);
        roomService.getRoomById(booking.getRoomId());

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
     * @throws HttpStatusException NOT_FOUND (почта/ауд. не найдена), UNAUTHORIZED (пароль не совпадает), FORBIDDEN (доступ запрещен)
     */
    public void updateBooking(Integer id, UpdateBookingDto booking, boolean isAdmin) throws HttpStatusException {
        
        auth(booking.getUserId(), booking.getEmail(), booking.getPassword(), isAdmin);
        roomService.getRoomById(booking.getRoomId());
        Booking bookingSql = getBookingById(id);

        if (booking.getBookingDate() != null) {
            bookingSql.setBookingDate(booking.getBookingDate());
        }
        if (booking.getBookingStart() != null) {
            bookingSql.setBookingStart(booking.getBookingStart());
        }
        if (booking.getBookingEnd() != null) {
            bookingSql.setBookingEnd(booking.getBookingEnd());
        }

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
