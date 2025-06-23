package centraluniversity.app.booking.models.rooms;

import javax.validation.constraints.NotBlank;
import javax.validation.constraints.NotNull;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
public class RoomCreateDto extends Room {

    @NotBlank
    private String name;

    @NotBlank
    private String description;

    @NotNull
    private Integer size;

    private String image;
}
