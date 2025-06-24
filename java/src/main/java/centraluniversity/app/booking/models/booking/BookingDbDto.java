package centraluniversity.app.booking.models.booking;

import java.time.LocalDate;
import java.time.LocalTime;

import com.fasterxml.jackson.annotation.JsonProperty;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Entity
@Table(name = "bookings")
@Setter
@Getter
@AllArgsConstructor
@NoArgsConstructor
public class BookingDbDto {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Integer id;

    @Column(nullable = false)
    @JsonProperty("user_id")
    private Integer userId;

    @JsonProperty("room_id")
    private Integer roomId;

    @Column(nullable = false)
    @JsonProperty("booking_date")
    private LocalDate bookingDate;

    @Column(nullable = false)
    @JsonProperty("booking_start")
    private LocalTime bookingStart;

    @Column(nullable = false)
    @JsonProperty("booking_end")
    private LocalTime bookingEnd;
    
}
