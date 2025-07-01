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
import centraluniversity.app.booking.repositories.BookingRepository;
import lombok.RequiredArgsConstructor;

@Service
@RequiredArgsConstructor
public class BookingService {

    private final RoomService roomService;
    private final BookingRepository bookingRepository;
    /**
     * Сохранение новой брони пользователем или администратором.
     * @param booking - информация о брони
     * @throws HttpStatusException CONFLICT (время брони занято), BAD_REQUEST (время start и end брони некорректны), NOT_FOUND (ауд. не найдена)
     */
    public void createBooking(BookingDbDto booking) throws HttpStatusException {
        
        // проверка существовании аудитории
        roomService.getRoomById(booking.getRoom().getId());

        // валидация времени и даты бронирования
        validateDateTime(booking.getBookingDate(), booking.getBookingStart(), booking.getBookingEnd());
        List<TimeBookingDto> bookingTimes = getTimesByDate(booking.getRoom().getId(), booking.getBookingDate()).getBookingTimes();
        validateTimesByOneDay(bookingTimes, booking.getBookingStart(), booking.getBookingEnd());

        bookingRepository.save(booking);
    }

    /**
     * Получение всех-всех-всех броней администратором.
     * @return пустой список или список с полной информацией о всех бронях
     */
    public List<BookingDbDto> getAllBookings() {
        return bookingRepository.findAll();
    }

    /**
     * Получение пользователем всех собственных броней
     * @return пустой список или список с полной информацией о всех бронях
     */
    public List<BookingDbDto> getAllBookingsByUserId(Integer userId) {
        return bookingRepository.findByUserId(userId);
    }

    /**
     * Получение информации о брони по его id.
     * @param id - id брони
     * @return полная информация о брони
     * @throws HttpStatusException NOT_FOUND (бронь не найдена)
     */
    public BookingDbDto getBookingById(Integer id) throws HttpStatusException {
        Optional<BookingDbDto> booking = bookingRepository.findById(id);
        if (booking.isEmpty()) {
            throw new HttpStatusException(HttpStatus.NOT_FOUND, String.format("Бронь с id = %d не существует", id));
        }
        return booking.get();
    }

    /**
     * Получение информации об бронях конкретной аудитории
     * @param id - id комнаты
     * @return пустой список или список 
     */
    public List<BookingDbDto> getAllBookingsByRoomId(Integer id) {
        return bookingRepository.findByRoomId(id);
    }

    /**
     * Обновление информации о бронировании пользователем или администратором.
     * @param bookingId - id брони
     * @param booking - полная информация о брони
     * @param isAdmin - отправлено админом или нет
     * @throws HttpStatusException CONFLICT (время брони занято), BAD_REQUEST (время start и end брони некорректны), NOT_FOUND (ауд./бронь не найдена), FORBIDDEN (при попытке изменять чужую запись не админом)
     */
    public void updateBooking(Integer bookingId, BookingDbDto booking) throws HttpStatusException {
        
        // проверка существовании аудитории
        roomService.getRoomById(bookingId);

        BookingDbDto bookingSql = getBookingById(bookingId);
        LocalTime oldBookingStart = bookingSql.getBookingStart();

        // замена полей в брони
        if (booking.getBookingDate() != null) {
            bookingSql.setBookingDate(booking.getBookingDate());
        }
        if (booking.getBookingStart() != null) {
            bookingSql.setBookingStart(booking.getBookingStart());
        }
        if (booking.getBookingEnd() != null) {
            bookingSql.setBookingEnd(booking.getBookingEnd());
        }

        // валидация дат и времен
        validateDateTime(bookingSql.getBookingDate(), bookingSql.getBookingStart(), bookingSql.getBookingEnd());
        List<TimeBookingDto> bookingTimes = getTimesByDate(booking.getRoom().getId(), booking.getBookingDate()).getBookingTimes();
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
     * @param id - id брони
     * @throws HttpStatusException NOT_FOUND (бронь не найдена)
     */
    public void deleteBooking(Integer id) throws HttpStatusException {
        getBookingById(id); // проверка существования
        bookingRepository.deleteById(id);
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
        List<BookingDbDto> bookings = bookingRepository.findByRoomIdAndBookingDateOrderByBookingStart(roomId, bookingDate);

        List<TimeBookingDto> bookingTimes = new ArrayList<>();
        for (int i = 0; i < bookings.size(); i++) {
            TimeBookingDto time = new TimeBookingDto();
            time.setBookingStart(bookings.get(i).getBookingStart());
            time.setBookingEnd(bookings.get(i).getBookingEnd());
            bookingTimes.add(time);
        }
        return new DateBookingDto(bookingDate, bookingTimes);
    }

    /**
     * Проверка на пересечения времени бронирования.
     * @param times - список всех интервалов за день
     * @param bookingStart
     * @param bookingEnd
     * @throws HttpStatusException CONFLICT (время брони занято)
     */
    public static void validateTimesByOneDay(List<TimeBookingDto> times, LocalTime bookingStart, LocalTime bookingEnd) throws HttpStatusException {
        DateTimeFormatter timeFormatter = DateTimeFormatter.ofPattern("HH:mm:ss");

        for (int i = 0; i < times.size(); i++) {
            TimeBookingDto time = times.get(i);
            boolean startIsBad = (bookingStart.isAfter(time.getBookingStart()) || bookingStart.equals(time.getBookingStart())) && (
                    bookingStart.isBefore(time.getBookingEnd()) || bookingStart.equals(time.getBookingEnd()));
            boolean endIsBad = (bookingEnd.isBefore(time.getBookingEnd()) || bookingEnd.equals(time.getBookingEnd())) && (
                    bookingEnd.isAfter(time.getBookingStart()) || bookingEnd.equals(time.getBookingStart()));
            boolean startAndEndAreBad = bookingStart.isBefore(time.getBookingStart()) && bookingEnd.isAfter(time.getBookingEnd());

            if (startIsBad || endIsBad || startAndEndAreBad) {
                throw new HttpStatusException(HttpStatus.CONFLICT, String.format("Уже существует бронь с %s по %s", 
                    time.getBookingStart().format(timeFormatter), time.getBookingEnd().format(timeFormatter)));
            }
        }
    }
}
