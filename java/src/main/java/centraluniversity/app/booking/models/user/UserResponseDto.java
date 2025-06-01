package centraluniversity.app.booking.models.user;

import centraluniversity.app.booking.pb.Role;
import lombok.*;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
public class UserResponseDto {
    
    private int id;
    private String email;
    private String firstName;
    private String secondName;
    private String patronymic;
    private String password;
    private Role role;
    
}
