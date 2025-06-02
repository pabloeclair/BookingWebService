package centraluniversity.app.booking.controllers;

import javax.validation.Valid;

import org.springframework.web.bind.annotation.*;

import centraluniversity.app.booking.models.user.CreateUserDto;
import centraluniversity.app.booking.models.user.GetUserAdminDto;
import centraluniversity.app.booking.models.user.UpdateUserDto;
import centraluniversity.app.booking.models.user.GetUserDto;
import centraluniversity.app.booking.pb.By;
import centraluniversity.app.booking.services.AdminUserService;
import lombok.RequiredArgsConstructor;

@RestController
@RequestMapping("/admin/users")
@RequiredArgsConstructor
public class AdminUserController {
    
    private final AdminUserService adminService;

    @PostMapping()
    public GetUserDto createUser(@Valid @RequestBody CreateUserDto user) {
        return adminService.createUser(user);
    }

    @GetMapping("/{by}/{sort_key}")
    public GetUserAdminDto getUser(@PathVariable("sort_key") String sortKey,
        @RequestParam(name = "email", required = true) String email,
        @RequestParam(name = "key", required = true) String key,
        @PathVariable String sortBy) {
        return adminService.getUser(email, key, By.valueOf(sortBy), sortKey);
    }

    @PutMapping("/{id}")
    public GetUserDto updateUser(@PathVariable("id") Integer id, @RequestBody @Valid UpdateUserDto user) {
        return adminService.updateUser(id, user);
    }

    @DeleteMapping("/{id}")
    public void deleteUser(@PathVariable("id") Integer id,
        @RequestParam(name = "email", required = true) String email,
        @RequestParam(name = "key", required = true) String key) {
        adminService.deleteUser(id, email, key);
    }
}
