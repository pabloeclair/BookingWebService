package centraluniversity.app.booking.controllers;

import java.util.List;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.*;

import centraluniversity.app.booking.models.booking.BookingCreateDto;
import centraluniversity.app.booking.models.booking.BookingDbDto;
import centraluniversity.app.booking.models.exception.HttpStatusException;
import centraluniversity.app.booking.models.rooms.RoomDbDto;
import centraluniversity.app.booking.models.user.UserDto;
import centraluniversity.app.booking.services.AuthService;
import centraluniversity.app.booking.services.BookingService;
import centraluniversity.app.booking.services.RoomService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;

@Tag(name = "Booking API", description = "Пользовательское управление бронями")
@RestController
@RequestMapping("/api/v1/bookings")
@RequiredArgsConstructor
public class BookingUserController {

    private final BookingService bookingService;
    private final RoomService roomService;
    private final AuthService authService;
    
    @Operation(summary = "Создание новой брони")
    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    public void createBooking(
            @RequestBody @Valid BookingCreateDto bookingDto,
            @RequestHeader("Authorization") String tokenString
    ) {
        UserDto user = authService.parseJwt(tokenString);
        if (!user.getId().equals(bookingDto.getUserId())) {
            throw new HttpStatusException(HttpStatus.FORBIDDEN, "бронировать разрешено только на себя");
        }
        RoomDbDto room = roomService.getRoomById(bookingDto.getRoomId());
        bookingService.createBooking(bookingDto.parseToDb(room));
    }

    @Operation(summary = "Получение списка всех бронь пользователя")
    @GetMapping("/user")
    public List<BookingDbDto> getBookingsByUserId(@RequestHeader("Authorization") String tokenString) {

        return bookingService.getAllBookingsByUserId(authService.parseJwt(tokenString).getId());
    }

    @Operation(summary = "Получение списка всех бронь комнаты")
    @GetMapping("/room/{id}")
    public List<BookingDbDto> getBookingsByRoomId(@PathVariable("id") Integer roomId) {

        return bookingService.getAllBookingsByRoomId(roomId);
    }

    @Operation(summary = "Обновление информации о брони")
    @PutMapping("/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void updateBooking(
            @PathVariable("id") Integer bookingId, 
            @RequestBody @Valid BookingDbDto booking,
            @RequestHeader("Authorization") String tokenString
    ) {
        auth(tokenString, bookingId);
        bookingService.updateBooking(bookingId, booking);
    }

    @Operation(summary = "Отмена брони")
    @DeleteMapping("/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void deleteBooking(
            @PathVariable("id") Integer bookingId,
            @RequestHeader("Authorization") String tokenString
    ) {
        auth(tokenString, bookingId);
        bookingService.deleteBooking(bookingId);
    }

    /**
     * Сверяет user_id из JWT и из брони, которую собирается изменить
     * @param tokenString
     * @param bookingId
     * @throws HttpStatusException NOT_FOUND (если ауд. не найден), FORBIDDEN (если не совпали id)
     */
    private void auth(String tokenString, Integer bookingId) throws HttpStatusException {
        UserDto user = authService.parseJwt(tokenString);
        BookingDbDto bookingSql = bookingService.getBookingById(bookingId);

        if (!user.getId().equals(bookingSql.getUserId())) {
            throw new HttpStatusException(HttpStatus.FORBIDDEN, "нельзя изменять чужую запись");
        }
    }
}
