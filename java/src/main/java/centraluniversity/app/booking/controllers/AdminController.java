package centraluniversity.app.booking.controllers;

import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import centraluniversity.app.booking.services.AdminService;
import lombok.RequiredArgsConstructor;

@RestController
@RequestMapping("/admin/users")
@RequiredArgsConstructor
public class AdminController {
    
    private final AdminService adminService;
}
