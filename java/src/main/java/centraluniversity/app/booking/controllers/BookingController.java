package centraluniversity.app.booking.controllers;

import java.util.List;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.*;

import centraluniversity.app.booking.models.booking.BookingDto;
import centraluniversity.app.booking.models.booking.CreateBookingDto;
import centraluniversity.app.booking.models.booking.UpdateBookingDto;
import centraluniversity.app.booking.services.BookingService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.parameters.RequestBody;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;

@Tag(name = "Booking API", description = "Пользовательское управление бронями")
@RestController
@RequestMapping("/bookings")
@RequiredArgsConstructor
public class BookingController {

    private final BookingService bookingService;
    
    @Operation(summary = "Создание новой брони")
    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    public void createBooking(@RequestBody @Valid CreateBookingDto bookingDto) throws Exception {
        bookingService.createBooking(bookingDto, false);
    }

    @Operation(summary = "Получение списка всех бронь пользователя")
    @GetMapping
    public List<BookingDto> getBookings(
            @RequestParam(name = "email", required = true) String email,
            @RequestParam(name = "key", required = true) String key) throws Exception {
        return bookingService.getAllBookingsByEmail(email, key);
    }

    @Operation(summary = "Обновление информации о брони")
    @PutMapping("/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void updateBooking(@PathVariable("id") Integer id, @RequestBody @Valid UpdateBookingDto booking) throws Exception {
        bookingService.updateBooking(id, booking, false);
    }

    @Operation(summary = "Отмена брони")
    @DeleteMapping("/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void deleteBooking(@PathVariable("id") Integer id,
            @RequestParam(name = "email", required = true) String email,
            @RequestParam(name = "key", required = true) String key) throws Exception {
        bookingService.deleteBooking(id, email, key, false);
    }
}
