package centraluniversity.app.booking.controllers;

import javax.validation.Valid;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.*;

import centraluniversity.app.booking.models.user.CreateUserDto;
import centraluniversity.app.booking.models.user.GetUserAdminDto;
import centraluniversity.app.booking.models.user.GetUserDto;
import centraluniversity.app.booking.models.user.UpdateRoleUserDto;
import centraluniversity.app.booking.models.user.UpdateUserDto;
import centraluniversity.app.booking.pb.By;
import centraluniversity.app.booking.services.AdminUserService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;

@Tag(name = "AdminUser API", description = "API управления пользователями")
@RestController
@RequestMapping("/admin/users")
@RequiredArgsConstructor
public class AdminUserController {
    
    private final AdminUserService adminService;

    @Operation(summary = "Регистрация нового пользователя")
    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    public void createUser(@Valid @RequestBody CreateUserDto user) throws Exception {
        adminService.createUser(user);
    }

    @Operation(summary = "Получение всех пользователей по ключу")
    @GetMapping("/{by}/{sort_key}")
    public GetUserAdminDto getUser(@PathVariable("sort_key") String sortKey,
            @RequestParam(name = "email", required = true) String email,
            @RequestParam(name = "key", required = true) String key,
            @PathVariable("by") String sortBy) throws Exception {
        return adminService.sortUser(email, key, By.valueOf(sortBy.toUpperCase()), sortKey);
    }

    @Operation(summary = "Редактирование пользователя")
    @PutMapping("/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
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

    @Operation(summary = "Изменение роли пользователя")
    @PutMapping("/role/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void updateRole(@PathVariable("id") Integer id,
            @RequestBody @Valid UpdateRoleUserDto user) {
        adminService.updateRole(id, user);
    }

    @Operation(summary = "Удаление пользователя")
    @DeleteMapping("/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void deleteUser(@PathVariable("id") Integer id,
            @RequestParam(name = "email", required = true) String email,
            @RequestParam(name = "key", required = true) String key) throws Exception {
        adminService.deleteUser(id, email, key);
    }
}
