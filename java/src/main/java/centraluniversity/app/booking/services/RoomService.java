package centraluniversity.app.booking.services;

import java.util.List;
import java.util.Optional;

import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;

import centraluniversity.app.booking.models.exception.HttpStatusException;
import centraluniversity.app.booking.models.rooms.CreateRoomDto;
import centraluniversity.app.booking.models.rooms.Room;
import centraluniversity.app.booking.models.rooms.UpdateRoomDto;
import centraluniversity.app.booking.models.user.GetUserDto;
import centraluniversity.app.booking.pb.Role;
import centraluniversity.app.booking.repositories.RoomRepository;
import lombok.RequiredArgsConstructor;

@Service
@RequiredArgsConstructor
public class RoomService {

    private final AuthService authService;
    private final RoomRepository roomRepository;
    
    public void createRoom(CreateRoomDto room) throws HttpStatusException {

        GetUserDto admin = authService.getUserByEmail(room.getAdminEmail(), room.getAdminPassword());
        if (admin.getRole() == Role.USER) {
            throw new HttpStatusException(HttpStatus.FORBIDDEN, "Создать комнату может только администратор");
        }

        Optional<Room> existingRoom = roomRepository.findByName(room.getName());
        if (existingRoom.isPresent()) {
            throw new HttpStatusException(HttpStatus.CONFLICT, "Комната с таким названием уже существует");
        }

        Room roomSql = new Room();
        roomSql.setName(room.getName());
        roomSql.setDescription(room.getDescription());
        roomSql.setSize(room.getSize());
        if (room.getImage() != null) {
            roomSql.setImage(room.getImage());
        }

        roomRepository.save(roomSql);
    }

    public List<Room> getAllRooms() {
        return roomRepository.findAll();
    }

    public List<Room> getRoomByName(String name) {
        return roomRepository.findByNameContaining(name);
    }

    public void updateRoom(Long roomId, UpdateRoomDto room) throws HttpStatusException {

        Optional<Room> optionalRoom = roomRepository.findById(roomId);
        if (optionalRoom.isEmpty()) {
            throw new HttpStatusException(HttpStatus.NOT_FOUND, "Комната не найдена");
        }

        Room roomSql = optionalRoom.get();
        
        if (room.getName() != null) {
            Optional<Room> existingRoom = roomRepository.findByName(room.getName());
            if (existingRoom.isPresent()) {
                throw new HttpStatusException(HttpStatus.CONFLICT, "Комната с таким названием уже существует");
            }
            roomSql.setName(room.getName());
        }
        
        if (room.getDescription() != null) {
            roomSql.setDescription(room.getDescription());
        }
        
        if (room.getSize() != null) {
            roomSql.setSize(room.getSize());
        }

        if (room.getImage() != null) {
            room.setImage(room.getImage());
        }

        roomRepository.save(roomSql);
    }

    public void deleteRoom(Long roomId, String adminEmail, String adminPassword) throws HttpStatusException {
        
        Optional<Room> optionalRoom = roomRepository.findById(roomId);
        if (optionalRoom.isEmpty()) {
            throw new HttpStatusException(HttpStatus.NOT_FOUND, "Комната не найдена");
        }

        GetUserDto admin = authService.getUserByEmail(adminEmail, adminPassword);
        if (admin.getRole() == Role.USER) {
            throw new HttpStatusException(HttpStatus.FORBIDDEN, "Создать комнату может только администратор");
        }

        roomRepository.deleteById(roomId);
    }
}
