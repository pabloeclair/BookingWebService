package centraluniversity.app.booking.models.admin;

import centraluniversity.app.booking.models.UserResponseDto;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class GetUserResponseDto {
    
    private UserResponseDto[] users;
}
