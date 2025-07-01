package centraluniversity.app.booking.models.booking;

import java.time.LocalDate;
import java.time.LocalTime;
import java.time.format.DateTimeFormatter;
import java.time.format.DateTimeParseException;

import org.springframework.http.HttpStatus;

import com.fasterxml.jackson.annotation.JsonProperty;

import centraluniversity.app.booking.models.exception.HttpStatusException;
import centraluniversity.app.booking.models.rooms.RoomDbDto;
import jakarta.validation.constraints.NotNull;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
public class BookingCreateDto {
    
    @NotNull
    @JsonProperty("user_id")
    private Integer userId;
    
    @NotNull
    @JsonProperty("room_id")
    private Integer roomId;

    @NotNull
    @JsonProperty("booking_date")
    private String bookingDate;

    @NotNull
    @JsonProperty("booking_start")
    private String bookingStart;

    @NotNull
    @JsonProperty("booking_end")
    private String bookingEnd;

    public BookingDbDto parseToDb(RoomDbDto room) throws HttpStatusException {

        DateTimeFormatter dateFormatter = DateTimeFormatter.ofPattern("dd.MM.yyyy");
        DateTimeFormatter timeFormatter = DateTimeFormatter.ofPattern("HH:mm:ss");
        LocalDate parsedBookingDate;
        LocalTime parsedBookingStart;
        LocalTime parsedBookingEnd;
        try {
            parsedBookingDate = LocalDate.parse(bookingDate, dateFormatter);
            parsedBookingStart = LocalTime.parse(bookingStart, timeFormatter);
            parsedBookingEnd = LocalTime.parse(bookingEnd, timeFormatter);
        } catch (DateTimeParseException e) {
            throw new HttpStatusException(HttpStatus.BAD_REQUEST, "неверный формат даты или времени; актуальный формат даты - dd.MM.yyyy, формат времени - HH:mm:ss");
        }
        return new BookingDbDto(null, userId, room, parsedBookingDate, parsedBookingStart, parsedBookingEnd);
    }
}
