package centraluniversity.app.booking.controllers;

import javax.validation.Valid;

import org.springframework.web.bind.annotation.*;

import centraluniversity.app.booking.models.user.SignupUserDto;
import centraluniversity.app.booking.models.user.GetUserDto;
import centraluniversity.app.booking.models.user.IdDto;
import centraluniversity.app.booking.services.AuthService;
import lombok.RequiredArgsConstructor;


@RestController
@RequestMapping("/users")
@RequiredArgsConstructor
public class AuthController {

    private final AuthService authService;

    @PostMapping()
    public IdDto signup(@Valid @RequestBody SignupUserDto userDto) throws Exception {
        return authService.createUser(userDto);
    }   

    @GetMapping()
    public GetUserDto login(
            @RequestParam(name = "email", required = true) String email,
            @RequestParam(name = "key", required = true) String key) throws Exception {
        return authService.getUserByEmail(email, key);
    }

}
