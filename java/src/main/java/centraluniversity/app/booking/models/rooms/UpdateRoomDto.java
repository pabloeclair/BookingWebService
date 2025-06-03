package centraluniversity.app.booking.models.rooms;

import javax.validation.constraints.Email;
import javax.validation.constraints.NotBlank;

import com.fasterxml.jackson.annotation.JsonProperty;

import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class UpdateRoomDto {
    
    private String name;

    private String description;

    private Integer size;

    private String image;

    @Email
    @NotBlank
    @JsonProperty("admin_email")
    private String adminEmail;

    @NotBlank
    @JsonProperty("admin_password")
    private String adminPassword;
}
