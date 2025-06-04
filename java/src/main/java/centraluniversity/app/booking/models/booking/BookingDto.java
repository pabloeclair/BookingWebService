package centraluniversity.app.booking.models.booking;

import java.time.LocalDate;
import java.time.LocalTime;

import com.fasterxml.jackson.annotation.JsonProperty;

import centraluniversity.app.booking.models.rooms.Room;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
@AllArgsConstructor
public class BookingDto {
    
    private Integer userId;
    
    @JsonProperty("room_id")
    private Integer roomId;

    @JsonProperty("booking_date")
    private LocalDate bookingDate;

    @JsonProperty("booking_start")
    private LocalTime bookingStart;

    @JsonProperty("booking_end")
    private LocalTime bookingEnd;

    private Room room;
}
