package centraluniversity.app.booking.controllers;

import javax.validation.Valid;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import centraluniversity.app.booking.models.auth.SignupUserDto;
import centraluniversity.app.booking.models.auth.UserResponseDto;
import centraluniversity.app.booking.services.AuthService;
import lombok.RequiredArgsConstructor;


@RestController
@RequestMapping("/users")
@RequiredArgsConstructor
public class AuthController {

    private final AuthService authService;

    @PostMapping()
    public UserResponseDto signup(@Valid @RequestBody SignupUserDto userDto) throws Exception {
        return authService.createUser(userDto);
    }   

    @GetMapping()
    public UserResponseDto getUserByEmail(@RequestParam String email) throws Exception {
        return authService.getUserByEmail(email);
    }

    @GetMapping("/{id}")
    public UserResponseDto getUserById(@PathVariable int id) throws Exception {
        return authService.getUserById(id);
    }
}
