package centraluniversity.app.booking.repositories;
import java.util.List;
import java.util.Optional;

import org.springframework.data.jpa.repository.JpaRepository;

import org.springframework.stereotype.Repository;

import centraluniversity.app.booking.models.rooms.RoomDbDto;

@Repository
public interface RoomRepository extends JpaRepository<RoomDbDto, Integer> {
    Optional<RoomDbDto> findByName(String name);
    List<RoomDbDto> findByNameContaining(String namePart);
}
