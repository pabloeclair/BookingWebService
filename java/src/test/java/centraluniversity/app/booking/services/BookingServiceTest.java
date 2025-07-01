package centraluniversity.app.booking.services;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

import java.time.LocalDate;
import java.time.LocalTime;
import java.time.format.DateTimeFormatter;
import java.util.ArrayList;
import java.util.List;

import org.junit.jupiter.api.Test;
import org.springframework.http.HttpStatus;

import centraluniversity.app.booking.models.booking.TimeBookingDto;
import centraluniversity.app.booking.models.exception.HttpStatusException;

public class BookingServiceTest {

    private final DateTimeFormatter dateFormatter = DateTimeFormatter.ofPattern("yyyy-MM-dd");
    
    // validateDateTime tests

    @Test
    public void testValidateDateTime_CorrectRequest() {

        // Array
        LocalDate bookingDate = LocalDate.parse("2044-04-04", dateFormatter);
        LocalTime bookingStart = LocalTime.of(8, 20);
        LocalTime bookingEnd = LocalTime.of(10, 20);

        // Action & Assert
        BookingService.validateDateTime(bookingDate, bookingStart, bookingEnd);
    }

    @Test
    public void testValidateDateTime_IncorrectRequest_IllegalBookingDate() {

        // Array
        LocalDate bookingDate = LocalDate.parse("2023-04-04", dateFormatter);
        LocalTime bookingStart = LocalTime.of(8, 20);
        LocalTime bookingEnd = LocalTime.of(10, 20);
        String expectedError = "Нельзя забронировать аудиторию ранее сегодняшнего дня";

        // Action & Assert
        HttpStatusException e = assertThrows(HttpStatusException.class,  
            () -> BookingService.validateDateTime(bookingDate, bookingStart, bookingEnd));
        assertEquals(HttpStatus.BAD_REQUEST, e.getStatus());
        assertEquals(expectedError, e.getMessage());
    }

    @Test
    public void testValidateDateTime_IncorrectRequest_IllegalBookingStart() {

        // Array
        LocalDate bookingDate = LocalDate.parse("2044-04-04", dateFormatter);
        LocalTime bookingStart = LocalTime.of(12, 20);
        LocalTime bookingEnd = LocalTime.of(10, 20);
        String expectedError = "Время начала должно быть строго раньше времени окончания";

        // Action & Assert
        HttpStatusException e = assertThrows(HttpStatusException.class,  
            () -> BookingService.validateDateTime(bookingDate, bookingStart, bookingEnd));
        assertEquals(HttpStatus.BAD_REQUEST, e.getStatus());
        assertEquals(expectedError, e.getMessage());
    }

    @Test
    public void testValidateDateTime_IncorrectRequest_BookingTimesEquals() {

        // Array
        LocalDate bookingDate = LocalDate.parse("2044-04-04", dateFormatter);
        LocalTime bookingStart = LocalTime.of(12, 20);
        LocalTime bookingEnd = LocalTime.of(12, 20);
        String expectedError = "Время начала должно быть строго раньше времени окончания";

        // Action & Assert
        HttpStatusException e = assertThrows(HttpStatusException.class,  
            () -> BookingService.validateDateTime(bookingDate, bookingStart, bookingEnd));
        assertEquals(HttpStatus.BAD_REQUEST, e.getStatus());
        assertEquals(expectedError, e.getMessage());
    }

    // validateTimesByOneDay tests

    private static List<TimeBookingDto> generateSampleTimeBookings() {
        List<TimeBookingDto> bookings = new ArrayList<>();

        TimeBookingDto booking2 = new TimeBookingDto();
        booking2.setBookingStart(LocalTime.of(9, 30));  
        booking2.setBookingEnd(LocalTime.of(10, 30));   
        bookings.add(booking2);

        TimeBookingDto booking3 = new TimeBookingDto();
        booking3.setBookingStart(LocalTime.of(11, 0));  
        booking3.setBookingEnd(LocalTime.of(12, 0));    
        bookings.add(booking3);

        TimeBookingDto booking4 = new TimeBookingDto();
        booking4.setBookingStart(LocalTime.of(14, 0));  
        booking4.setBookingEnd(LocalTime.of(15, 0));    
        bookings.add(booking4);

        return bookings;
    }

    @Test
    public void testValidateTimesByOneDay_CorrectRequest_WithArray() {

        // Array
        List<TimeBookingDto> array = generateSampleTimeBookings();
        LocalTime bookingStart = LocalTime.of(16, 20);
        LocalTime bookingEnd = LocalTime.of(18, 20);

        // Action & Assert
        BookingService.validateTimesByOneDay(array, bookingStart, bookingEnd);
    }

    @Test
    public void testValidateTimesByOneDay_CorrectRequest_WithoutArray() {

        // Array
        List<TimeBookingDto> array = new ArrayList<>();
        LocalTime bookingStart = LocalTime.of(16, 20);
        LocalTime bookingEnd = LocalTime.of(18, 20);

        // Action & Assert
        BookingService.validateTimesByOneDay(array, bookingStart, bookingEnd);
    }

    @Test 
    public void testValidateTimesByOneDay_IncorrectRequest_IllegalStartTime1() {

        // Array
        List<TimeBookingDto> array = generateSampleTimeBookings();
        LocalTime bookingStart = LocalTime.of(11, 20);
        LocalTime bookingEnd = LocalTime.of(13, 20);
        String expectedError = "Уже существует бронь с 11:00:00 по 12:00:00";

        // Action & Assert
        HttpStatusException e = assertThrows(HttpStatusException.class, () -> BookingService.validateTimesByOneDay(array, bookingStart, bookingEnd));
        assertEquals(HttpStatus.CONFLICT, e.getStatus());
        assertEquals(expectedError, e.getMessage());
    }

    @Test 
    public void testValidateTimesByOneDay_IncorrectRequest_IllegalStartTime2() {

        // Array
        List<TimeBookingDto> array = generateSampleTimeBookings();
        LocalTime bookingStart = LocalTime.of(12, 0);
        LocalTime bookingEnd = LocalTime.of(13, 20);
        String expectedError = "Уже существует бронь с 11:00:00 по 12:00:00";

        // Action & Assert
        HttpStatusException e = assertThrows(HttpStatusException.class, () -> BookingService.validateTimesByOneDay(array, bookingStart, bookingEnd));
        assertEquals(HttpStatus.CONFLICT, e.getStatus());
        assertEquals(expectedError, e.getMessage());
    }
    
    @Test 
    public void testValidateTimesByOneDay_IncorrectRequest_IllegalEndTime1() {

        // Array
        List<TimeBookingDto> array = generateSampleTimeBookings();
        LocalTime bookingStart = LocalTime.of(8, 20);
        LocalTime bookingEnd = LocalTime.of(9, 40);
        String expectedError = "Уже существует бронь с 09:30:00 по 10:30:00";

        // Action & Assert
        HttpStatusException e = assertThrows(HttpStatusException.class, () -> BookingService.validateTimesByOneDay(array, bookingStart, bookingEnd));
        assertEquals(HttpStatus.CONFLICT, e.getStatus());
        assertEquals(expectedError, e.getMessage());
    }

    @Test 
    public void testValidateTimesByOneDay_IncorrectRequest_IllegalEndTime2() {

        // Array
        List<TimeBookingDto> array = generateSampleTimeBookings();
        LocalTime bookingStart = LocalTime.of(8, 20);
        LocalTime bookingEnd = LocalTime.of(9, 30);
        String expectedError = "Уже существует бронь с 09:30:00 по 10:30:00";

        // Action & Assert
        HttpStatusException e = assertThrows(HttpStatusException.class, () -> BookingService.validateTimesByOneDay(array, bookingStart, bookingEnd));
        assertEquals(HttpStatus.CONFLICT, e.getStatus());
        assertEquals(expectedError, e.getMessage());
    }

    @Test 
    public void testValidateTimesByOneDay_IncorrectRequest_IllegalStartAndEndTime() {

        // Array
        List<TimeBookingDto> array = generateSampleTimeBookings();
        LocalTime bookingStart = LocalTime.of(8, 20);
        LocalTime bookingEnd = LocalTime.of(10, 40);
        String expectedError = "Уже существует бронь с 09:30:00 по 10:30:00";

        // Action & Assert
        HttpStatusException e = assertThrows(HttpStatusException.class, () -> BookingService.validateTimesByOneDay(array, bookingStart, bookingEnd));
        assertEquals(HttpStatus.CONFLICT, e.getStatus());
        assertEquals(expectedError, e.getMessage());
    }
}
