package centraluniversity.app.booking.repositories;

import java.time.LocalDate;
import java.util.List;

import org.springframework.data.jpa.repository.JpaRepository;

import centraluniversity.app.booking.models.booking.Booking;

public interface BookingRepository extends JpaRepository<Booking, Integer> {
    List<Booking> findByUserId(Integer userId);
    void deleteByRoomId(Integer roomId);
    void deleteByUserId(Integer userId);
    List<Booking> findByRoomIdAndBookingDateOrderByBookingStart(Integer roomId, LocalDate bookingDate);
}
