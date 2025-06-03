package centraluniversity.app.booking.controllers;

import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import centraluniversity.app.booking.services.BookingService;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;

@Tag(name = "Booking API", description = "Пользовательское управление бронями")
@RestController
@RequestMapping("/books")
@RequiredArgsConstructor
public class BookingController {

    private final BookingService bookingService;
    
}
