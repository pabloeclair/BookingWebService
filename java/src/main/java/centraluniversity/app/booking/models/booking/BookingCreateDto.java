package centraluniversity.app.booking.models.booking;

import java.time.LocalDate;
import java.time.LocalTime;

import com.fasterxml.jackson.annotation.JsonProperty;

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
    private Integer userId;
    
    @NotNull
    @JsonProperty("room_id")
    private Integer roomId;

    @NotNull
    @JsonProperty("booking_date")
    private LocalDate bookingDate;

    @NotNull
    @JsonProperty("booking_start")
    private LocalTime bookingStart;

    @NotNull
    @JsonProperty("booking_end")
    private LocalTime bookingEnd;

    @NotNull
    private RoomDbDto room;

    public BookingDbDto parseToDb()  {
        return new BookingDbDto(roomId, userId, roomId, bookingDate, bookingStart, bookingEnd);
    }
}
