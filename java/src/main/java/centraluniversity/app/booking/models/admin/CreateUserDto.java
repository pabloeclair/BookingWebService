package centraluniversity.app.booking.models.admin;

import javax.validation.constraints.Email;
import javax.validation.constraints.NotBlank;

import com.fasterxml.jackson.annotation.JsonProperty;

import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class CreateUserDto {

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

    @Email
    @NotBlank
    @JsonProperty("admin_email")
    private String adminEmail;

    @NotBlank
    @JsonProperty("admin_password")
    private String adminPassword;
}
