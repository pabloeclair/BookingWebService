package centraluniversity.app.booking.models.booking;

import java.time.LocalTime;

import com.fasterxml.jackson.annotation.JsonProperty;

import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class TimeBookingDto {
    
    @JsonProperty("booking_start")
    private LocalTime bookingStart;

    @JsonProperty("booking_end")
    private LocalTime bookingEnd;
}
