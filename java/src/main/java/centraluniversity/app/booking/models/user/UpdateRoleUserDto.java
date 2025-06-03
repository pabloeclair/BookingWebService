package centraluniversity.app.booking.models.user;

import javax.validation.constraints.Email;
import javax.validation.constraints.NotBlank;

import com.fasterxml.jackson.annotation.JsonProperty;

import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class UpdateRoleUserDto {
    
    @NotBlank
    @Email
    @JsonProperty("admin_email")
    private String adminEmail;

    @NotBlank
    @JsonProperty("admin_password")
    private String adminPassword;

    @NotBlank
    @JsonProperty("old_role")
    private String oldRole;

    @NotBlank 
    @JsonProperty("new_role")
    private String newRole;
}
