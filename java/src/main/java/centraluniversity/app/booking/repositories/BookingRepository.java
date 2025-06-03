package centraluniversity.app.booking.repositories;

import org.springframework.data.jpa.repository.JpaRepository;

import centraluniversity.app.booking.models.booking.Booking;

public interface BookingRepository extends JpaRepository<Booking, Long> {
}
