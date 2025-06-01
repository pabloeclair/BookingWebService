package centraluniversity.app.booking.models.user;

import javax.validation.constraints.Email;
import javax.validation.constraints.NotBlank;

import com.fasterxml.jackson.annotation.JsonProperty;

import centraluniversity.app.booking.pb.Role;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class SignupUserDto {
    
    @Email
    @NotBlank
    private String email;

    @NotBlank
    @JsonProperty("first_name")
    private String firstName;

    @NotBlank
    @JsonProperty("second_name")
    private String secondName;

    private String patronymic;

    @NotBlank
    private String password;

    private Role role;
}
