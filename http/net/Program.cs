using System.Text.Json;

var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();
app.MapGet("/", () => "Hello from .NET minimal API");
app.MapGet("/json", () => JsonSerializer.Serialize("Hello from .NET minimal API JSON"));
app.Run();
