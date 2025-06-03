package centraluniversity.app.booking.models.rooms;

import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class RoomDto {
    private Integer id;
    private String name;
    private String description;
    private int size;
    private String image;
}
