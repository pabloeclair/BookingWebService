package centraluniversity.app.booking.controllers;

import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import centraluniversity.app.booking.services.AdminBookingService;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;

@Tag(name = "AdminBooking API", description = "Административное управление бронями")
@RestController
@RequestMapping("/admin/books")
@RequiredArgsConstructor
public class AdminBookingController {

    private final AdminBookingService bookingService;
    
}
