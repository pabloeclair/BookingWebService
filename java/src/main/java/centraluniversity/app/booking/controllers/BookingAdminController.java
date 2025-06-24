package centraluniversity.app.booking.controllers;

import java.util.List;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.*;

import centraluniversity.app.booking.models.booking.BookingCreateDto;
import centraluniversity.app.booking.models.booking.BookingDbDto;
import centraluniversity.app.booking.models.exception.HttpStatusException;
import centraluniversity.app.booking.models.user.Role;
import centraluniversity.app.booking.models.user.UserDto;
import centraluniversity.app.booking.services.AuthService;
import centraluniversity.app.booking.services.BookingService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;

@Tag(name = "AdminBooking API", description = "Административное управление бронями")
@RestController
@RequestMapping("/api/v1/admin/booking")
@RequiredArgsConstructor
public class BookingAdminController {

    private final BookingService bookingService;
    private final AuthService authService;
    
    @Operation(summary = "Создание новой брони")
    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    public void createBooking(
            @RequestBody @Valid BookingCreateDto bookingDto,
            @RequestHeader("Authorization") String tokenString
    ) {
        auth(tokenString);
        bookingService.createBooking(bookingDto.parseToDb());
    }

    @Operation(summary = "Получение всех доступных броней")
    @GetMapping
    public List<BookingDbDto> getAllBookings(@RequestHeader("Authorization") String tokenString) {
        auth(tokenString);
        return bookingService.getAllBookings();
    }

    @Operation(summary = "Обновление информации о брони")
    @PutMapping("/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void updateBooking(
            @PathVariable("id") Integer id, 
            @RequestBody @Valid BookingDbDto booking,
            @RequestHeader("Authorization") String tokenString
    ) {
        auth(tokenString);
        bookingService.updateBooking(id, booking);
    }

    @Operation(summary = "Отмена брони")
    @DeleteMapping("/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void deleteBooking(
            @PathVariable("id") Integer id,
            @RequestHeader("Authorization") String tokenString
    ) {
        auth(tokenString);
        bookingService.deleteBooking(id);
    }

    private void auth(String tokenString) throws HttpStatusException {
        UserDto user = authService.parseJwt(tokenString);
        if (!user.getRole().equals(Role.ADMIN) && !user.getRole().equals(Role.MAIN_ADMIN)) {
            throw new HttpStatusException(HttpStatus.FORBIDDEN, "доступ запрещен");
        }
    }
}
