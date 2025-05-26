package centraluniversity.app.booking.models.exception;

import lombok.Getter;
import org.springframework.http.HttpStatus;

@Getter
public class HttpStatusException extends RuntimeException {

    private final HttpStatus status;

    public HttpStatusException(HttpStatus status, String message) {
        super(message);
        this.status = status;
    }
}

