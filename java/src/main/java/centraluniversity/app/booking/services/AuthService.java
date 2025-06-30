package centraluniversity.app.booking.services;

import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;

import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;

import com.fasterxml.jackson.databind.ObjectMapper;

import centraluniversity.app.booking.models.exception.ErrorDto;
import centraluniversity.app.booking.models.exception.HttpStatusException;
import centraluniversity.app.booking.models.user.UserDto;
import lombok.NoArgsConstructor;

@Service
@NoArgsConstructor
public class AuthService {
    private String url = "http://user-api:7070/api/v1/user";

    public UserDto parseJwt(String tokenString) {
        HttpResponse<String> response;
        try {
            HttpClient client = HttpClient.newHttpClient();
            HttpRequest request = HttpRequest.newBuilder()
                    .uri(URI.create(url))
                    .header("Authorization", tokenString)
                    .GET()
                    .build();
            response = client.send(request, HttpResponse.BodyHandlers.ofString());
        } catch (Exception e) {
            throw new HttpStatusException(HttpStatus.INTERNAL_SERVER_ERROR, e.getMessage());
        }
        return handleResponse(response);
    }

    private UserDto handleResponse(HttpResponse<String> response) {
        if (response.statusCode() == 200) {
            try {
                ObjectMapper objectMapper = new ObjectMapper();
                return objectMapper.readValue(response.body(), UserDto.class);
            } catch (Exception e) {
                throw new HttpStatusException(HttpStatus.INTERNAL_SERVER_ERROR, "ошибка парсинга ответа: " + e.getMessage());
            }
        } else {
            ErrorDto errorDto;
            try {
                ObjectMapper objectMapper = new ObjectMapper();
                errorDto = objectMapper.readValue(response.body(), ErrorDto.class);
            } catch (Exception e) {
                throw new HttpStatusException(HttpStatus.INTERNAL_SERVER_ERROR, "ошибка парсинга ошибки: " + e.getMessage());
            }
            throw new HttpStatusException(HttpStatus.valueOf(response.statusCode()), errorDto.getErrorMessage());
        }
    }
    
} 
