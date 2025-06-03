package centraluniversity.app.booking.models.user;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
@AllArgsConstructor
public class GetUserAdminDto {
    
    private GetUserDto[] users;
}
