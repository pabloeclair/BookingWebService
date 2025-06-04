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
    
    public Room getRoomById(Integer id) {
        Optional<Room> room = roomRepository.findById(id);
        if (room.isEmpty()) {
            throw new HttpStatusException(HttpStatus.NOT_FOUND, String.format("Аудитории с id = %d не существует", id));
        }
        return room.get();
    }

    private void auth(String email, String password) throws HttpStatusException {
        GetUserDto admin = authService.getUserByEmail(email, password);
        if (admin.getRole() == Role.USER) {
            throw new HttpStatusException(HttpStatus.FORBIDDEN, "Создать аудиторию может только администратор");
        }
    }

    public void createRoom(CreateRoomDto room) throws HttpStatusException {

        auth(room.getAdminEmail(), room.getAdminPassword());
        Optional<Room> existingRoom = roomRepository.findByName(room.getName());
        if (existingRoom.isPresent()) {
            throw new HttpStatusException(HttpStatus.CONFLICT, String.format("Аудитория с названием '%s' уже существует", room.getName()));
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

    public void updateRoom(Integer roomId, UpdateRoomDto room) throws HttpStatusException {

        auth(room.getAdminEmail(), room.getAdminPassword());
        Room roomSql = getRoomById(roomId);
        
        if (room.getName() != null) {
            Optional<Room> existingRoom = roomRepository.findByName(room.getName());
            if (existingRoom.isPresent()) {
                throw new HttpStatusException(HttpStatus.CONFLICT, String.format("Аудитория с названием '%s' уже существует", room.getName()));
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

    public void deleteRoom(Integer roomId, String adminEmail, String adminPassword) throws HttpStatusException {   
        auth(adminEmail, adminPassword);
        getRoomById(roomId);
        roomRepository.deleteById(roomId);
    }
}
