package centraluniversity.app.booking.controllers;

import java.util.List;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.*;

import centraluniversity.app.booking.models.exception.HttpStatusException;
import centraluniversity.app.booking.models.rooms.RoomDbDto;
import centraluniversity.app.booking.models.rooms.RoomCreateDto;
import centraluniversity.app.booking.models.user.Role;
import centraluniversity.app.booking.models.user.UserDto;
import centraluniversity.app.booking.services.AuthService;
import centraluniversity.app.booking.services.RoomService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;

@Tag(name = "Room API", description = "API управления аудиториями для бронирования")
@RestController
@RequestMapping("/api/v1")
@RequiredArgsConstructor
public class RoomController {

    private final RoomService roomService;
    private final AuthService authService;
    
    @Operation(summary = "Создание новой комнаты")
    @PostMapping("/admin/rooms")
    @ResponseStatus(HttpStatus.CREATED)
    public void createRoom(
            @RequestBody @Valid RoomCreateDto room, 
            @RequestHeader("Authorization") String tokenString
    ) {
        auth(tokenString);
        roomService.createRoom(room);
    }

    @Operation(summary = "Получение всех аудиторий по названию")
    @GetMapping("/rooms")
    public List<RoomDbDto> getRoomsByName(@RequestParam(value = "name", required = false) String partName) {
        if (partName == null) {
            return roomService.getAllRooms();
        }
        return roomService.getRoomByName(partName);
    }  

    @Operation(summary = "Получение аудитории по id")
    @GetMapping("/rooms/{id}")
    public RoomDbDto getRoomById(@PathVariable("id") Integer id) {
        return roomService.getRoomById(id);
    }

    @Operation(summary = "Изменение данных об аудитории")
    @PutMapping("/admin/rooms/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void updateRoom(
            @PathVariable("id") Integer id, 
            @RequestBody RoomDbDto room,
            @RequestHeader("Authorization") String tokenString
    ) {
        auth(tokenString);
        roomService.updateRoom(id, room);
    }

    @Operation(summary = "Удаление аудитории")
    @DeleteMapping("/admin/rooms/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void deleteRoom(
            @PathVariable("id") Integer id, 
            @RequestHeader("Authorization") String tokenString
    ) {
        auth(tokenString);
        roomService.deleteRoom(id);
    }

    /**
     * Проверяет роль пользователя по JWT токену
     * @param tokenString
     * @throws HttpStatusException UNAUTHORIZED (если что-то не так с JWT токеном), FORBIDDEN (если роль не подходит)
     */
    private void auth(String tokenString) throws HttpStatusException {
        UserDto user = authService.parseJwt(tokenString);
        if (!user.getRole().equals(Role.ADMIN) && !user.getRole().equals(Role.MAIN_ADMIN)) {
            throw new HttpStatusException(HttpStatus.FORBIDDEN, "доступ запрещен");
        }
    }
}
