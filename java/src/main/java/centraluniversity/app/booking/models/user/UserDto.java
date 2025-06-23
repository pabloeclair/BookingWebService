package centraluniversity.app.booking.models.user;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import lombok.*;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
@JsonIgnoreProperties(ignoreUnknown = true)
public class UserDto {
    
    private Integer id;
    private String email;
    
    @JsonProperty("first_name")
    private String firstName;

    @JsonProperty("second_name")
    private String secondName;
    
    private String patronymic;
    private Role role;
    
}
