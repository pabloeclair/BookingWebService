package centraluniversity.app.booking.services;

import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;

import centraluniversity.app.booking.models.exception.HttpStatusException;
import centraluniversity.app.booking.models.user.SignupUserDto;
import centraluniversity.app.booking.models.user.UserResponseDto;
import centraluniversity.app.booking.pb.AuthenticationGrpc;
import centraluniversity.app.booking.pb.Email;
import centraluniversity.app.booking.pb.SignupRequest;
import centraluniversity.app.booking.pb.UserResponse;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.grpc.Status;
import io.grpc.StatusRuntimeException;
import jakarta.annotation.PostConstruct;
import jakarta.annotation.PreDestroy;

@Service
public class AuthService {
    
    private ManagedChannel channel;
    private AuthenticationGrpc.AuthenticationBlockingStub stub;

    @PostConstruct
    public void connectToServer() {
        this.channel = ManagedChannelBuilder.forAddress("auth", 7001)
            .usePlaintext()
            .build();
        this.stub = AuthenticationGrpc.newBlockingStub(channel);
    }

    @PreDestroy
    public void shutdown() {
        if (this.channel != null) {
            this.channel.shutdown();
        }
    }

    /**
     * Create new service user
     * @param user - signup user information
     * @return UserResponseDto - full user information
     * @throws Exception
     */
    public UserResponseDto createUser(SignupUserDto user) {

        SignupRequest req;
        if (user.getPatronymic() == null || user.getPatronymic().isEmpty()) {
            req = SignupRequest.newBuilder()
                .setEmail(user.getEmail())
                .setFirstName(user.getFirstName())
                .setSecondName(user.getSecondName())
                .setPassword(user.getPassword())
                .build();
        } else {
            req = SignupRequest.newBuilder()
                .setEmail(user.getEmail())
                .setFirstName(user.getFirstName())
                .setSecondName(user.getSecondName())
                .setPatronymic(user.getPatronymic())
                .setPassword(user.getPassword())
                .build();
        }

        UserResponse res;
        try {
            res = this.stub.signupUser(req);
        } catch (StatusRuntimeException e) {
            Status status = e.getStatus();
            if (status.getCode() == Status.Code.ALREADY_EXISTS) {
                throw new HttpStatusException(HttpStatus.BAD_REQUEST, e.getMessage());
            } 
            throw new HttpStatusException(HttpStatus.INTERNAL_SERVER_ERROR, e.getMessage());
        } 
        
        return new UserResponseDto(
            res.getId(),
            res.getEmail(),
            res.getFirstName(),
            res.getSecondName(),
            res.getPatronymic(),
            res.getPassword(),
            res.getRole()
        );  
    }

    /**
     * Get user by email
     * @param email
     * @return UserResponseDto - full user information
     * @throws Exception
     */
    public UserResponseDto getUserByEmail(String email, String password) throws Exception {

        Email req = Email.newBuilder().setEmail(email).setPassword(password).build();
        UserResponse res;
        try {
            res = this.stub.loginUser(req);
        } catch (StatusRuntimeException e) {
            Status status = e.getStatus();
            switch (status.getCode()) {
                case NOT_FOUND:
                    throw new HttpStatusException(HttpStatus.NOT_FOUND, e.getMessage());
                case UNAUTHENTICATED:
                    throw new HttpStatusException(HttpStatus.UNAUTHORIZED, e.getMessage());
                default:
                    throw new HttpStatusException(HttpStatus.INTERNAL_SERVER_ERROR, e.getMessage());
            }
        }
        
        return new UserResponseDto(
            res.getId(),
            res.getEmail(),
            res.getFirstName(),
            res.getSecondName(),
            res.getPatronymic(),
            res.getPassword(),
            res.getRole()
        ); 
    }

}
