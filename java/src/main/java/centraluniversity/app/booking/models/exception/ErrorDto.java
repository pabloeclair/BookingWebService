package centraluniversity.app.booking.models.exception;

import com.fasterxml.jackson.annotation.JsonProperty;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
@AllArgsConstructor
public class ErrorDto {

    @JsonProperty("status_code")
    private String statusCode;

    @JsonProperty("error_message")
    private String errorMessage;
}

