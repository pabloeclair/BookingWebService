package centraluniversity.app.booking.controllers;

import centraluniversity.app.booking.models.exception.ErrorDto;
import centraluniversity.app.booking.models.exception.HttpStatusException;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;

@ControllerAdvice
public class CommonExceptionController {

    @ExceptionHandler(MethodArgumentNotValidException.class)
    public ResponseEntity<ErrorDto> methodArgumentNotValidException(MethodArgumentNotValidException e) {
        return ResponseEntity.status(e.getStatusCode())
                .body(new ErrorDto(e.getStatusCode().toString(), null, e.getMessage()));
    }

    @ExceptionHandler(HttpStatusException.class)
    public ResponseEntity<ErrorDto> handleHttpStatusException(HttpStatusException e) {
        return ResponseEntity.status(e.getStatus())
                .body(new ErrorDto(e.getStatus().toString(), null, e.getMessage()));
    }

    @ExceptionHandler(Exception.class)
    public ResponseEntity<ErrorDto> handleException(Exception e) {
        return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                .body(new ErrorDto("500 INTERNAL_SERVER_ERROR", e.getClass().getName(), e.getMessage()));
    }
}
