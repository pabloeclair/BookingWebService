package centraluniversity.app.booking.models.rooms;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
public class RoomCreateDto {

    @NotBlank
    private String name;

    @NotBlank
    private String description;

    @NotNull
    private Integer size;

    private String image;

    public Room parseToDb() {
        Room room = new Room(this.name, this.description, this.size);
        if (this.image != null) {
            room.setImage(this.image);
        }
        return room;
    }
}
