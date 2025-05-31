package centraluniversity.app.booking.models.admin;

import javax.validation.constraints.Email;
import javax.validation.constraints.NotBlank;

import com.fasterxml.jackson.annotation.JsonProperty;

import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class DeleteUserDto {
    
    @NotBlank
    private Integer id;

    @Email
    @NotBlank
    @JsonProperty("admin_email")
    private String adminEmail;

    @NotBlank
    @JsonProperty("admin_password")
    private String adminPassword;
}
