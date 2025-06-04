package centraluniversity.app.booking.controllers;

import javax.validation.Valid;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.*;

import centraluniversity.app.booking.models.user.SignupUserDto;
import centraluniversity.app.booking.models.user.GetUserDto;
import centraluniversity.app.booking.models.user.IdDto;
import centraluniversity.app.booking.services.AuthService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;

@Tag(name = "Auth API", description = "API авторизации пользователей")
@RestController
@RequestMapping("/users")
@RequiredArgsConstructor
public class AuthController {

    private final AuthService authService;

    @Operation(summary = "Регистрация пользователя")
    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    public IdDto signup(@Valid @RequestBody SignupUserDto userDto) throws Exception {
        return authService.createUser(userDto);
    }   

    @Operation(summary = "Получение пользователя по почте и паролю")
    @GetMapping
    public GetUserDto login(
            @RequestParam(name = "email", required = true) String email,
            @RequestParam(name = "key", required = true) String key) throws Exception {
        return authService.getUserByEmail(email, key);
    }

}
