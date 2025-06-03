package centraluniversity.app.booking.controllers;

import java.util.List;

import javax.validation.Valid;

import org.springframework.web.bind.annotation.*;

import centraluniversity.app.booking.models.rooms.CreateRoomDto;
import centraluniversity.app.booking.models.rooms.Room;
import centraluniversity.app.booking.models.rooms.UpdateRoomDto;
import centraluniversity.app.booking.services.RoomService;
import lombok.RequiredArgsConstructor;

@RestController
@RequiredArgsConstructor
public class RoomController {

    private final RoomService roomService;
    
    @PostMapping("/admin/rooms")
    public void createRoom(CreateRoomDto room) {
        roomService.createRoom(room);
    }

    @GetMapping("/rooms")
    public List<Room> getAllRooms() {
        return roomService.getAllRooms();
    }

    @GetMapping("/rooms/{part_name}")
    public List<Room> getRoomsByName(@PathVariable("part_name") String partName) {
        return roomService.getRoomByName(partName);
    }    

    @PutMapping("/admin/rooms/{id}")
    public void updateRoom(@PathVariable("id") Long id, @RequestBody @Valid UpdateRoomDto room) {
        roomService.updateRoom(id, room);
    }

    @DeleteMapping("/admin/rooms/{id}")
    public void deleteRoom(@PathVariable("id") Long id, 
            @RequestParam(name = "email", required = true) String email,
            @RequestParam(name = "key", required = true) String key) {
        roomService.deleteRoom(id, email, key);
    }
}
