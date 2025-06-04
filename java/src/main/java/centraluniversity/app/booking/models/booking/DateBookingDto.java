package centraluniversity.app.booking.models.booking;

import java.time.LocalDate;
import java.util.List;

import com.fasterxml.jackson.annotation.JsonProperty;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
@AllArgsConstructor
public class DateBookingDto {
    
    @JsonProperty("booking_date")
    private LocalDate bookingDate;

    @JsonProperty("booking_times")
    private List<TimeBookingDto> bookingTimes;
}
