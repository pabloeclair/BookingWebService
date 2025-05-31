package centraluniversity.app.booking.models.admin;

import javax.validation.constraints.Email;
import javax.validation.constraints.NotBlank;

import com.fasterxml.jackson.annotation.JsonProperty;

import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class GetUserRequestDto {
    
    @NotBlank
    @JsonProperty("sort_by")
    private String sortBy;

    @NotBlank
    @JsonProperty("sort_key")
    private String sortKey;

    @Email
    @NotBlank
    @JsonProperty("admin_email")
    private String adminEmail;

    @NotBlank
    @JsonProperty("admin_password")
    private String adminPassword;
}
