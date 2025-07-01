package centraluniversity.app.booking.repositories;

import java.time.LocalDate;
import java.util.List;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import centraluniversity.app.booking.models.booking.BookingDbDto;

@Repository
public interface BookingRepository extends JpaRepository<BookingDbDto, Integer> {
    List<BookingDbDto> findByUserId(Integer userId);
    void deleteByRoomId(Integer roomId);
    void deleteByUserId(Integer userId);
    List<BookingDbDto> findByRoomIdAndBookingDateOrderByBookingStart(Integer roomId, LocalDate bookingDate);
    List<BookingDbDto> findByRoomId(Integer roomId);
}
