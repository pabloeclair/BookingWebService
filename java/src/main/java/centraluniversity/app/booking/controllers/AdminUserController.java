package centraluniversity.app.booking.controllers;

import javax.validation.Valid;

import org.springframework.web.bind.annotation.*;

import centraluniversity.app.booking.models.user.CreateUserDto;
import centraluniversity.app.booking.models.user.GetUserAdminDto;
import centraluniversity.app.booking.models.user.GetUserDto;
import centraluniversity.app.booking.models.user.UpdateUserDto;
import centraluniversity.app.booking.pb.By;
import centraluniversity.app.booking.services.AdminUserService;
import lombok.RequiredArgsConstructor;

@RestController
@RequestMapping("/admin/users")
@RequiredArgsConstructor
public class AdminUserController {
    
    private final AdminUserService adminService;

    @PostMapping()
    public void createUser(@Valid @RequestBody CreateUserDto user) throws Exception {
        adminService.createUser(user);
    }

    @GetMapping("/{by}/{sort_key}")
    public GetUserAdminDto getUser(@PathVariable("sort_key") String sortKey,
            @RequestParam(name = "email", required = true) String email,
            @RequestParam(name = "key", required = true) String key,
            @PathVariable("by") String sortBy) throws Exception {
        return adminService.sortUser(email, key, By.valueOf(sortBy.toUpperCase()), sortKey);
    }

    @PutMapping("/{id}")
    public void updateUser(@PathVariable("id") Integer id,
            @RequestBody @Valid UpdateUserDto user) throws Exception {
        GetUserDto oldUser = adminService.getUserById(id, user.getAdminEmail(), user.getAdminPassword());

        System.out.println(user.getFirstName());
        if (user.getEmail() == null) {
            user.setEmail(oldUser.getEmail());
        }
        if (user.getFirstName() == null) {
            user.setFirstName(oldUser.getFirstName());
        }
        if (user.getSecondName() == null) {
            user.setSecondName(oldUser.getSecondName());
        }
        if (user.getPatronymic() == null && oldUser.getPatronymic() != null) {
            user.setPatronymic(oldUser.getPatronymic());
        }
        if (user.getPassword() == null) {
            user.setPassword(oldUser.getPassword());
        }
        if (user.getRole() == null) {
            user.setRole(oldUser.getRole().toString());
        }
        user.setId(id);

        adminService.updateUser(user);
    }

    @DeleteMapping("/{id}")
    public void deleteUser(@PathVariable("id") Integer id,
            @RequestParam(name = "email", required = true) String email,
            @RequestParam(name = "key", required = true) String key) throws Exception {
        adminService.deleteUser(id, email, key);
    }
}
