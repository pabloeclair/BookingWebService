package centraluniversity.app.booking.models.user;

import com.fasterxml.jackson.annotation.JsonProperty;

import centraluniversity.app.booking.pb.Role;
import lombok.*;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
public class GetUserDto {
    
    private Integer id;
    private String email;
    
    @JsonProperty("first_name")
    private String firstName;

    @JsonProperty("second_name")
    private String secondName;
    
    private String patronymic;
    private String password;
    private Role role;
    
}
