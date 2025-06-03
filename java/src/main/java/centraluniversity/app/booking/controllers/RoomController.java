package centraluniversity.app.booking.controllers;

import java.util.List;

import javax.validation.Valid;

import org.springframework.web.bind.annotation.*;

import centraluniversity.app.booking.models.rooms.CreateRoomDto;
import centraluniversity.app.booking.models.rooms.Room;
import centraluniversity.app.booking.models.rooms.UpdateRoomDto;
import centraluniversity.app.booking.services.RoomService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;

@Tag(name = "Room API", description = "API управления аудиториями для бронирования")
@RestController
@RequiredArgsConstructor
public class RoomController {

    private final RoomService roomService;
    
    @Operation(summary = "Создание новой комнаты")
    @PostMapping("/admin/rooms")
    public void createRoom(CreateRoomDto room) {
        roomService.createRoom(room);
    }

    @Operation(summary = "Получение всех аудиторий")
    @GetMapping("/rooms")
    public List<Room> getAllRooms() {
        return roomService.getAllRooms();
    }

    @Operation(summary = "Получение аудиторий по названию")
    @GetMapping("/rooms/{part_name}")
    public List<Room> getRoomsByName(@PathVariable("part_name") String partName) {
        return roomService.getRoomByName(partName);
    }    

    @Operation(summary = "Изменение данных об аудитории")
    @PutMapping("/admin/rooms/{id}")
    public void updateRoom(@PathVariable("id") Long id, @RequestBody @Valid UpdateRoomDto room) {
        roomService.updateRoom(id, room);
    }

    @Operation(summary = "Удаление аудитории")
    @DeleteMapping("/admin/rooms/{id}")
    public void deleteRoom(@PathVariable("id") Long id, 
            @RequestParam(name = "email", required = true) String email,
            @RequestParam(name = "key", required = true) String key) {
        roomService.deleteRoom(id, email, key);
    }
}
