package centraluniversity.app.booking.controllers;

import centraluniversity.app.booking.models.exception.ErrorDto;
import centraluniversity.app.booking.models.exception.HttpStatusException;

import java.util.List;
import java.util.stream.Collectors;

import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.http.converter.HttpMessageNotReadableException;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.MissingRequestHeaderException;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.servlet.resource.NoResourceFoundException;

@ControllerAdvice
public class CommonExceptionController {

    // Обработчик HttpStatusException
    @ExceptionHandler(HttpStatusException.class)
    public ResponseEntity<ErrorDto> handleHttpStatusException(HttpStatusException e) {
        return ResponseEntity.status(e.getStatus())
                .body(new ErrorDto(e.getStatus().toString(), e.getMessage()));
    }

    // Обработчик отсутствия заголовка
    @ExceptionHandler(MissingRequestHeaderException.class)
    public ResponseEntity<ErrorDto> handleMissingRequestHeaderException(MissingRequestHeaderException e) {
        return ResponseEntity.status(HttpStatus.BAD_REQUEST)
                .body(new ErrorDto(HttpStatus.BAD_REQUEST.toString(), e.getMessage()));
    }

    // Обработчик не пройденной валидации
    @ExceptionHandler(MethodArgumentNotValidException.class)
    public ResponseEntity<ErrorDto> handleMethodArgumentNotValidException(MethodArgumentNotValidException e) {
        List<String> errorFields = e.getBindingResult().getFieldErrors()
                .stream()
                .map(error -> error.getField() + ": " + error.getDefaultMessage())
                .collect(Collectors.toList());

        String errorMessage = String.join(", ", errorFields);
        return ResponseEntity.status(HttpStatus.BAD_REQUEST)
                .body(new ErrorDto(HttpStatus.BAD_REQUEST.toString(), errorMessage));
    }

    // Обработчик ошибки невозможности десериализации
    @ExceptionHandler(HttpMessageNotReadableException.class)
    public ResponseEntity<ErrorDto> handleMessageNotReadableException(HttpMessageNotReadableException e) {
        return ResponseEntity.status(HttpStatus.BAD_REQUEST)
                .body(new ErrorDto(HttpStatus.BAD_REQUEST.toString(), e.getMessage()));
    }

    // Обработчик не найденного адресного пути
    @ExceptionHandler(NoResourceFoundException.class)
    public ResponseEntity<ErrorDto> handleNoResourceFoundException(NoResourceFoundException e) {
        return ResponseEntity.status(HttpStatus.NOT_FOUND)
                .body(new ErrorDto(HttpStatus.NOT_FOUND.toString(), e.getMessage()));
    }

    // Обработчик остальных ошибок
    @ExceptionHandler(Exception.class)
    public ResponseEntity<ErrorDto> handleException(Exception e) {
        return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                .body(new ErrorDto(HttpStatus.INTERNAL_SERVER_ERROR.toString(), e.getClass().getName() + ": " + e.getMessage()));
    }
}
