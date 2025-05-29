package centraluniversity.app.booking.models.auth;

import javax.validation.constraints.Email;
import javax.validation.constraints.NotBlank;

import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class LoginUserDto {

    @Email
    @NotBlank
    private String email;

    @NotBlank
    private String password;
    
}
