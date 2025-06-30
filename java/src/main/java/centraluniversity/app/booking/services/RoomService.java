package centraluniversity.app.booking.services;

import java.util.List;
import java.util.Optional;

import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;

import centraluniversity.app.booking.models.exception.HttpStatusException;
import centraluniversity.app.booking.models.rooms.RoomDbDto;
import centraluniversity.app.booking.models.rooms.RoomCreateDto;
import centraluniversity.app.booking.repositories.BookingRepository;
import centraluniversity.app.booking.repositories.RoomRepository;
import lombok.RequiredArgsConstructor;

@Service
@RequiredArgsConstructor
public class RoomService {

    private final RoomRepository roomRepository;
    private final BookingRepository bookingRepository;
    
    /**
     * Получение информации об аудитории по id.
     * @param id аудитории
     * @return Room - подробная информация об аудитории
     * @throws HttpStatusException NOT_FOUND (не найдена аудитория)
     */
    public RoomDbDto getRoomById(Integer id) throws HttpStatusException {
        Optional<RoomDbDto> room = roomRepository.findById(id);
        if (room.isEmpty()) {
            throw new HttpStatusException(HttpStatus.NOT_FOUND, String.format("Аудитории с id = %d не существует", id));
        }
        return room.get();
    }

    /**
     * Сохранение новой аудитории.
     * @param room - полная информация о сохраняемой комнате
     * @throws HttpStatusException CONFLICT (указанное название ауд. уже существует)
     */
    public void createRoom(RoomCreateDto room) throws HttpStatusException {

        Optional<RoomDbDto> existingRoom = roomRepository.findByName(room.getName());
        if (existingRoom.isPresent()) {
            throw new HttpStatusException(HttpStatus.CONFLICT, String.format("аудитория с названием '%s' уже существует", room.getName()));
        }
        
        roomRepository.save(room.parseToDb());
    }

    /**
     * Получение информации обо всех аудиториях.
     * @return пустой список, если ничего не найдено, или список всех аудиторий с полной информацией о них
     */
    public List<RoomDbDto> getAllRooms() {
        return roomRepository.findAll();
    }

    /**
     * Получение всех аудиторий, схожих по указанному названию.
     * @param name - полное или частичное название аудитории
     * @return пустой список, если ничего не найдено, или список всех найденных аудиторий с полной информацией о них 
     */
    public List<RoomDbDto> getRoomByName(String name) {
        return roomRepository.findByNameContaining(name);
    }

    /**
     * Обновление информации об указанной аудитории.
     * @param roomId - id обновляемой аудитории
     * @param room - новая информация об аудитории
     * @throws HttpStatusException FORBIDDEN (отказано в доступе), CONFLICT (указанное название ауд. уже существует)
     */
    public void updateRoom(Integer roomId, RoomDbDto room) throws HttpStatusException {

        RoomDbDto roomSql = getRoomById(roomId);
        
        if (room.getName() != null) {
            Optional<RoomDbDto> existingRoom = roomRepository.findByName(room.getName());
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
            roomSql.setImage(room.getImage());
        }

        roomRepository.save(roomSql);
    }

    /**
     * Удаление любой информации об аудитории в таблицах rooms и bookings.
     * @param roomId
     * @throws HttpStatusException FORBIDDEN (отказано в доступе), NOT_FOUND (не найдена аудитория)
     */
    public void deleteRoom(Integer roomId) throws HttpStatusException {   

        // Проверка, что комната существует
        getRoomById(roomId);

        // Action
        roomRepository.deleteById(roomId);
        bookingRepository.deleteByRoomId(roomId);
    }
}
