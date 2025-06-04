package centraluniversity.app.booking.models.booking;

import java.time.LocalDate;
import java.time.LocalTime;

import com.fasterxml.jackson.annotation.JsonProperty;

import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotBlank;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class CreateBookingDto {
    
    @NotBlank
    @JsonProperty("user_id")
    private Integer userId;
    
    @NotBlank
    @JsonProperty("room_id")
    private Integer roomId;

    @NotBlank
    @JsonProperty("booking_date")
    private LocalDate bookingDate;

    @NotBlank 
    @JsonProperty("booking_start")
    private LocalTime bookingStart;

    @NotBlank
    @JsonProperty("booking_end")
    private LocalTime bookingEnd;

    @Email
    @NotBlank
    private String email;

    @NotBlank
    private String password;
}
