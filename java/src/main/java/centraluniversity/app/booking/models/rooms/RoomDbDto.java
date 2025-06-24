package centraluniversity.app.booking.models.rooms;

import jakarta.persistence.*;
import lombok.*;

@Entity 
@Table(name = "rooms") 
@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
public class RoomDbDto {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Integer id;

    @Column(nullable = false, unique = true)
    private String name;

    @Column(nullable = false)
    private String description;

    @Column(nullable = false)
    private Integer size;

    @Column
    private String image;

    RoomDbDto(String name, String description, Integer size) {
        this.name = name;
        this.description = description;
        this.size = size;
    }
    
}
