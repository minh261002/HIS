package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/minhtran/his/internal/middleware"
	"github.com/minhtran/his/internal/pkg/jwt"
	"github.com/minhtran/his/internal/pkg/response"
)

// SetupRoutes configures all application routes
func SetupRoutes(
	r *gin.Engine,
	authHandler *AuthHandler,
	userHandler *UserHandler,
	patientHandler *PatientHandler,
	allergyHandler *PatientAllergyHandler,
	historyHandler *PatientMedicalHistoryHandler,
	appointmentHandler *AppointmentHandler,
	visitHandler *VisitHandler,
	icd10Handler *ICD10CodeHandler,
	diagnosisHandler *DiagnosisHandler,
	medicationHandler *MedicationHandler,
	prescriptionHandler *PrescriptionHandler,
	labTestTemplateHandler *LabTestTemplateHandler,
	labTestRequestHandler *LabTestRequestHandler,
	imagingTemplateHandler *ImagingTemplateHandler,
	imagingRequestHandler *ImagingRequestHandler,
	bedHandler *BedHandler,
	admissionHandler *AdmissionHandler,
	inventoryHandler *InventoryHandler,
	dispensingHandler *DispensingHandler,
	invoiceHandler *InvoiceHandler,
	paymentHandler *PaymentHandler,
	insuranceClaimHandler *InsuranceClaimHandler,
	departmentHandler *DepartmentHandler,
	medicalServiceHandler *MedicalServiceHandler,
	auditLogHandler *AuditLogHandler,
	jwtManager *jwt.Manager,
	rbacMiddleware *middleware.RBACMiddleware,
	allowedOrigins []string,
) {
	// Apply global middleware
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.CORSMiddleware(allowedOrigins))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, "Server is healthy", map[string]string{
			"status": "ok",
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Public auth routes (only login and refresh)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(jwtManager))
		{
			// Auth protected routes
			protected.GET("/auth/profile", authHandler.GetProfile)

			// User management routes (admin only)
			users := protected.Group("/users")
			users.Use(rbacMiddleware.RequirePermission("users.manage"))
			{
				users.POST("", userHandler.CreateUser)
				users.GET("", userHandler.ListUsers)
				users.GET("/:id", userHandler.GetUser)
				users.PUT("/:id", userHandler.UpdateUser)
				users.DELETE("/:id", userHandler.DeleteUser)
				users.POST("/:id/roles", userHandler.AssignRoles)
			}

			// Patient management routes
			patients := protected.Group("/patients")
			{
				// Statistics (requires view permission)
				patients.GET("/stats", rbacMiddleware.RequirePermission("patients.view"), patientHandler.GetPatientStats)

				// Search (requires view permission)
				patients.GET("/search", rbacMiddleware.RequirePermission("patients.view"), patientHandler.SearchPatients)

				// Get by code (requires view permission)
				patients.GET("/code/:code", rbacMiddleware.RequirePermission("patients.view"), patientHandler.GetPatientByCode)

				// Register patient (requires create permission)
				patients.POST("", rbacMiddleware.RequirePermission("patients.create"), patientHandler.RegisterPatient)

				// List patients (requires view permission)
				patients.GET("", rbacMiddleware.RequirePermission("patients.view"), patientHandler.ListPatients)

				// Get patient details (requires view permission)
				patients.GET("/:id", rbacMiddleware.RequirePermission("patients.view"), patientHandler.GetPatient)

				// Update patient (requires update permission)
				patients.PUT("/:id", rbacMiddleware.RequirePermission("patients.update"), patientHandler.UpdatePatient)

				// Delete patient (requires delete permission - admin only)
				patients.DELETE("/:id", rbacMiddleware.RequirePermission("patients.delete"), patientHandler.DeletePatient)

				// Patient allergies sub-routes
				patients.POST("/:id/allergies", rbacMiddleware.RequirePermission("patients.create"), allergyHandler.AddAllergy)
				patients.GET("/:id/allergies", rbacMiddleware.RequirePermission("patients.view"), allergyHandler.GetPatientAllergies)
				patients.GET("/:id/allergies/active", rbacMiddleware.RequirePermission("patients.view"), allergyHandler.GetActiveAllergies)

				// Patient medical history sub-routes
				patients.POST("/:id/medical-history", rbacMiddleware.RequirePermission("patients.create"), historyHandler.AddMedicalHistory)
				patients.GET("/:id/medical-history", rbacMiddleware.RequirePermission("patients.view"), historyHandler.GetPatientHistory)
				patients.GET("/:id/medical-history/active", rbacMiddleware.RequirePermission("patients.view"), historyHandler.GetActiveConditions)
			}

			// Allergy routes (standalone)
			allergies := protected.Group("/allergies")
			{
				allergies.GET("/:allergyId", rbacMiddleware.RequirePermission("patients.view"), allergyHandler.GetAllergy)
				allergies.PUT("/:allergyId", rbacMiddleware.RequirePermission("patients.update"), allergyHandler.UpdateAllergy)
				allergies.DELETE("/:allergyId", rbacMiddleware.RequirePermission("patients.delete"), allergyHandler.DeleteAllergy)
			}

			// Medical history routes (standalone)
			medicalHistory := protected.Group("/medical-history")
			{
				medicalHistory.GET("/:historyId", rbacMiddleware.RequirePermission("patients.view"), historyHandler.GetMedicalHistory)
				medicalHistory.PUT("/:historyId", rbacMiddleware.RequirePermission("patients.update"), historyHandler.UpdateMedicalHistory)
				medicalHistory.DELETE("/:historyId", rbacMiddleware.RequirePermission("patients.delete"), historyHandler.DeleteMedicalHistory)
			}

			// Appointment routes
			appointments := protected.Group("/appointments")
			{
				// List and search
				appointments.GET("", rbacMiddleware.RequirePermission("appointments.view"), appointmentHandler.ListAppointments)
				appointments.GET("/upcoming", rbacMiddleware.RequirePermission("appointments.view"), appointmentHandler.GetUpcomingAppointments)
				appointments.GET("/code/:code", rbacMiddleware.RequirePermission("appointments.view"), appointmentHandler.GetAppointmentByCode)

				// Create
				appointments.POST("", rbacMiddleware.RequirePermission("appointments.create"), appointmentHandler.ScheduleAppointment)

				// Get details
				appointments.GET("/:id", rbacMiddleware.RequirePermission("appointments.view"), appointmentHandler.GetAppointment)

				// Update/Reschedule
				appointments.PUT("/:id", rbacMiddleware.RequirePermission("appointments.update"), appointmentHandler.RescheduleAppointment)

				// Status transitions
				appointments.POST("/:id/cancel", rbacMiddleware.RequirePermission("appointments.cancel"), appointmentHandler.CancelAppointment)
				appointments.POST("/:id/confirm", rbacMiddleware.RequirePermission("appointments.manage"), appointmentHandler.ConfirmAppointment)
				appointments.POST("/:id/start", rbacMiddleware.RequirePermission("appointments.manage"), appointmentHandler.StartAppointment)
				appointments.POST("/:id/complete", rbacMiddleware.RequirePermission("appointments.manage"), appointmentHandler.CompleteAppointment)
				appointments.POST("/:id/no-show", rbacMiddleware.RequirePermission("appointments.manage"), appointmentHandler.MarkNoShow)
			}

			// Patient appointments sub-routes
			protected.GET("/patients/:id/appointments", rbacMiddleware.RequirePermission("appointments.view"), appointmentHandler.GetPatientAppointments)

			// Doctor schedule routes
			protected.GET("/doctors/:id/schedule", rbacMiddleware.RequirePermission("appointments.view"), appointmentHandler.GetDoctorSchedule)
			protected.GET("/doctors/:id/available-slots", rbacMiddleware.RequirePermission("appointments.view"), appointmentHandler.GetAvailableTimeSlots)

			// Visit routes
			visits := protected.Group("/visits")
			{
				// List and search
				visits.GET("", rbacMiddleware.RequirePermission("visits.view"), visitHandler.ListVisits)
				visits.GET("/code/:code", rbacMiddleware.RequirePermission("visits.view"), visitHandler.GetVisitByCode)

				// Create
				visits.POST("", rbacMiddleware.RequirePermission("visits.create"), visitHandler.CreateVisit)

				// Get details
				visits.GET("/:id", rbacMiddleware.RequirePermission("visits.view"), visitHandler.GetVisit)

				// Update
				visits.PUT("/:id", rbacMiddleware.RequirePermission("visits.update"), visitHandler.UpdateVisit)

				// Status transitions
				visits.POST("/:id/complete", rbacMiddleware.RequirePermission("visits.complete"), visitHandler.CompleteVisit)
				visits.POST("/:id/cancel", rbacMiddleware.RequirePermission("visits.delete"), visitHandler.CancelVisit)
			}

			// Patient visits sub-routes
			protected.GET("/patients/:id/visits", rbacMiddleware.RequirePermission("visits.view"), visitHandler.GetPatientVisits)

			// Doctor visits routes
			protected.GET("/doctors/:id/visits", rbacMiddleware.RequirePermission("visits.view"), visitHandler.GetDoctorVisits)

			// ICD-10 code routes
			icd10 := protected.Group("/icd10-codes")
			{
				icd10.GET("/search", rbacMiddleware.RequirePermission("diagnoses.view"), icd10Handler.SearchICD10Codes)
				icd10.GET("/:code", rbacMiddleware.RequirePermission("diagnoses.view"), icd10Handler.GetICD10CodeByCode)
				icd10.GET("/category/:category", rbacMiddleware.RequirePermission("diagnoses.view"), icd10Handler.GetICD10CodesByCategory)
			}

			// Diagnosis routes
			diagnoses := protected.Group("/diagnoses")
			{
				diagnoses.POST("", rbacMiddleware.RequirePermission("diagnoses.create"), diagnosisHandler.AddDiagnosis)
				diagnoses.GET("/:id", rbacMiddleware.RequirePermission("diagnoses.view"), diagnosisHandler.GetDiagnosis)
				diagnoses.PUT("/:id", rbacMiddleware.RequirePermission("diagnoses.update"), diagnosisHandler.UpdateDiagnosis)
				diagnoses.DELETE("/:id", rbacMiddleware.RequirePermission("diagnoses.delete"), diagnosisHandler.DeleteDiagnosis)
			}

			// Visit/Patient diagnosis sub-routes
			protected.GET("/visits/:id/diagnoses", rbacMiddleware.RequirePermission("diagnoses.view"), diagnosisHandler.GetVisitDiagnoses)
			protected.GET("/patients/:id/diagnoses", rbacMiddleware.RequirePermission("diagnoses.view"), diagnosisHandler.GetPatientDiagnoses)

			// Medication routes
			medications := protected.Group("/medications")
			{
				medications.GET("/search", rbacMiddleware.RequirePermission("prescriptions.view"), medicationHandler.SearchMedications)
				medications.GET("/:id", rbacMiddleware.RequirePermission("prescriptions.view"), medicationHandler.GetMedication)
			}

			// Prescription routes
			prescriptions := protected.Group("/prescriptions")
			{
				prescriptions.POST("", rbacMiddleware.RequirePermission("prescriptions.create"), prescriptionHandler.CreatePrescription)
				prescriptions.GET("/:id", rbacMiddleware.RequirePermission("prescriptions.view"), prescriptionHandler.GetPrescription)
				prescriptions.GET("/code/:code", rbacMiddleware.RequirePermission("prescriptions.view"), prescriptionHandler.GetPrescriptionByCode)
				prescriptions.PUT("/:id", rbacMiddleware.RequirePermission("prescriptions.update"), prescriptionHandler.UpdatePrescription)
				prescriptions.POST("/:id/dispense", rbacMiddleware.RequirePermission("prescriptions.dispense"), prescriptionHandler.DispensePrescription)
				prescriptions.POST("/:id/complete", rbacMiddleware.RequirePermission("prescriptions.dispense"), prescriptionHandler.CompletePrescription)
				prescriptions.POST("/:id/cancel", rbacMiddleware.RequirePermission("prescriptions.delete"), prescriptionHandler.CancelPrescription)
			}

			// Visit/Patient prescription sub-routes
			protected.GET("/visits/:id/prescriptions", rbacMiddleware.RequirePermission("prescriptions.view"), prescriptionHandler.GetVisitPrescriptions)
			protected.GET("/patients/:id/prescriptions", rbacMiddleware.RequirePermission("prescriptions.view"), prescriptionHandler.GetPatientPrescriptions)

			// Lab test template routes
			labTestTemplates := protected.Group("/lab-test-templates")
			{
				labTestTemplates.GET("/search", rbacMiddleware.RequirePermission("lab_tests.view"), labTestTemplateHandler.SearchTemplates)
				labTestTemplates.GET("/:code", rbacMiddleware.RequirePermission("lab_tests.view"), labTestTemplateHandler.GetTemplateByCode)
				labTestTemplates.GET("/category/:category", rbacMiddleware.RequirePermission("lab_tests.view"), labTestTemplateHandler.GetTemplatesByCategory)
			}

			// Lab test request routes
			labTestRequests := protected.Group("/lab-test-requests")
			{
				labTestRequests.POST("", rbacMiddleware.RequirePermission("lab_tests.create"), labTestRequestHandler.CreateLabTestRequest)
				labTestRequests.GET("/:id", rbacMiddleware.RequirePermission("lab_tests.view"), labTestRequestHandler.GetLabTestRequest)
				labTestRequests.GET("/code/:code", rbacMiddleware.RequirePermission("lab_tests.view"), labTestRequestHandler.GetLabTestRequestByCode)
				labTestRequests.POST("/:id/collect-sample", rbacMiddleware.RequirePermission("lab_tests.enter_results"), labTestRequestHandler.CollectSample)
				labTestRequests.POST("/:id/start-processing", rbacMiddleware.RequirePermission("lab_tests.enter_results"), labTestRequestHandler.StartProcessing)
				labTestRequests.POST("/:id/complete", rbacMiddleware.RequirePermission("lab_tests.enter_results"), labTestRequestHandler.CompleteTest)
				labTestRequests.POST("/:id/cancel", rbacMiddleware.RequirePermission("lab_tests.delete"), labTestRequestHandler.CancelTest)
				labTestRequests.POST("/:id/results", rbacMiddleware.RequirePermission("lab_tests.enter_results"), labTestRequestHandler.EnterResults)
			}

			// Visit/Patient lab test sub-routes
			protected.GET("/visits/:id/lab-tests", rbacMiddleware.RequirePermission("lab_tests.view"), labTestRequestHandler.GetVisitLabTests)
			protected.GET("/patients/:id/lab-tests", rbacMiddleware.RequirePermission("lab_tests.view"), labTestRequestHandler.GetPatientLabTests)

			// Imaging template routes
			imagingTemplates := protected.Group("/imaging-templates")
			{
				imagingTemplates.GET("/search", rbacMiddleware.RequirePermission("imaging.view"), imagingTemplateHandler.SearchTemplates)
				imagingTemplates.GET("/:code", rbacMiddleware.RequirePermission("imaging.view"), imagingTemplateHandler.GetTemplateByCode)
				imagingTemplates.GET("/modality/:modality", rbacMiddleware.RequirePermission("imaging.view"), imagingTemplateHandler.GetTemplatesByModality)
			}

			// Imaging request routes
			imagingRequests := protected.Group("/imaging-requests")
			{
				imagingRequests.POST("", rbacMiddleware.RequirePermission("imaging.create"), imagingRequestHandler.CreateImagingRequest)
				imagingRequests.GET("/:id", rbacMiddleware.RequirePermission("imaging.view"), imagingRequestHandler.GetImagingRequest)
				imagingRequests.GET("/code/:code", rbacMiddleware.RequirePermission("imaging.view"), imagingRequestHandler.GetImagingRequestByCode)
				imagingRequests.POST("/:id/schedule", rbacMiddleware.RequirePermission("imaging.update"), imagingRequestHandler.ScheduleImaging)
				imagingRequests.POST("/:id/start", rbacMiddleware.RequirePermission("imaging.report"), imagingRequestHandler.StartImaging)
				imagingRequests.POST("/:id/complete", rbacMiddleware.RequirePermission("imaging.report"), imagingRequestHandler.CompleteImaging)
				imagingRequests.POST("/:id/cancel", rbacMiddleware.RequirePermission("imaging.delete"), imagingRequestHandler.CancelImaging)
				imagingRequests.POST("/:id/result", rbacMiddleware.RequirePermission("imaging.report"), imagingRequestHandler.CreateOrUpdateResult)
			}

			// Visit/Patient imaging sub-routes
			protected.GET("/visits/:id/imaging-requests", rbacMiddleware.RequirePermission("imaging.view"), imagingRequestHandler.GetVisitImagingRequests)
			protected.GET("/patients/:id/imaging-requests", rbacMiddleware.RequirePermission("imaging.view"), imagingRequestHandler.GetPatientImagingRequests)

			// Bed routes
			beds := protected.Group("/beds")
			{
				beds.GET("/available", rbacMiddleware.RequirePermission("beds.view"), bedHandler.GetAvailableBeds)
				beds.GET("/:number", rbacMiddleware.RequirePermission("beds.view"), bedHandler.GetBedByNumber)
			}

			// Admission routes
			admissions := protected.Group("/admissions")
			{
				admissions.POST("", rbacMiddleware.RequirePermission("admissions.create"), admissionHandler.CreateAdmission)
				admissions.GET("/:id", rbacMiddleware.RequirePermission("admissions.view"), admissionHandler.GetAdmission)
				admissions.GET("/code/:code", rbacMiddleware.RequirePermission("admissions.view"), admissionHandler.GetAdmissionByCode)
				admissions.POST("/:id/discharge", rbacMiddleware.RequirePermission("admissions.discharge"), admissionHandler.DischargeAdmission)
				admissions.POST("/:id/transfer-bed", rbacMiddleware.RequirePermission("beds.manage"), admissionHandler.TransferBed)
				admissions.GET("/active", rbacMiddleware.RequirePermission("admissions.view"), admissionHandler.GetActiveAdmissions)
				admissions.POST("/:id/nursing-notes", rbacMiddleware.RequirePermission("nursing_notes.create"), admissionHandler.CreateNursingNote)
				admissions.GET("/:id/nursing-notes", rbacMiddleware.RequirePermission("nursing_notes.view"), admissionHandler.GetAdmissionNursingNotes)
			}

			// Patient admission history sub-route
			protected.GET("/patients/:id/admissions", rbacMiddleware.RequirePermission("admissions.view"), admissionHandler.GetPatientAdmissions)

			// Inventory routes
			inventory := protected.Group("/inventory")
			{
				inventory.POST("", rbacMiddleware.RequirePermission("inventory.manage"), inventoryHandler.AddStock)
				inventory.GET("/medication/:id", rbacMiddleware.RequirePermission("inventory.view"), inventoryHandler.GetMedicationStock)
				inventory.GET("/low-stock", rbacMiddleware.RequirePermission("inventory.view"), inventoryHandler.GetLowStockAlerts)
				inventory.GET("/expiring-soon", rbacMiddleware.RequirePermission("inventory.view"), inventoryHandler.GetExpiringSoonAlerts)
			}

			// Dispensing routes
			dispensing := protected.Group("/dispensing")
			{
				dispensing.POST("", rbacMiddleware.RequirePermission("dispensing.dispense"), dispensingHandler.DispensePrescription)
			}

			// Prescription/Patient dispensing sub-routes
			protected.GET("/prescriptions/:id/dispensing", rbacMiddleware.RequirePermission("dispensing.view"), dispensingHandler.GetPrescriptionDispensingRecords)
			protected.GET("/patients/:id/dispensing", rbacMiddleware.RequirePermission("dispensing.view"), dispensingHandler.GetPatientDispensingHistory)

			// Invoice routes
			invoices := protected.Group("/invoices")
			{
				invoices.POST("", rbacMiddleware.RequirePermission("invoices.create"), invoiceHandler.CreateInvoice)
				invoices.GET("/:id", rbacMiddleware.RequirePermission("invoices.view"), invoiceHandler.GetInvoice)
				invoices.GET("/code/:code", rbacMiddleware.RequirePermission("invoices.view"), invoiceHandler.GetInvoiceByCode)
			}

			// Payment routes
			payments := protected.Group("/payments")
			{
				payments.POST("", rbacMiddleware.RequirePermission("payments.process"), paymentHandler.CreatePayment)
			}

			// Insurance claim routes
			insuranceClaims := protected.Group("/insurance-claims")
			{
				insuranceClaims.POST("", rbacMiddleware.RequirePermission("insurance_claims.manage"), insuranceClaimHandler.CreateInsuranceClaim)
				insuranceClaims.POST("/:id/approve", rbacMiddleware.RequirePermission("insurance_claims.manage"), insuranceClaimHandler.ApproveClaim)
				insuranceClaims.POST("/:id/reject", rbacMiddleware.RequirePermission("insurance_claims.manage"), insuranceClaimHandler.RejectClaim)
			}

			// Patient/Invoice sub-routes
			protected.GET("/patients/:id/invoices", rbacMiddleware.RequirePermission("invoices.view"), invoiceHandler.GetPatientInvoices)
			protected.GET("/invoices/:id/payments", rbacMiddleware.RequirePermission("payments.view"), paymentHandler.GetInvoicePayments)
			protected.GET("/invoices/:id/insurance-claims", rbacMiddleware.RequirePermission("insurance_claims.view"), insuranceClaimHandler.GetInvoiceClaims)

			// System Module Routes
			system := protected.Group("/system")
			{
				// Departments
				depts := system.Group("/departments")
				{
					depts.POST("", rbacMiddleware.RequirePermission("departments.create"), departmentHandler.CreateDepartment)
					depts.GET("", rbacMiddleware.RequirePermission("departments.view"), departmentHandler.ListDepartments)
					depts.GET("/:id", rbacMiddleware.RequirePermission("departments.view"), departmentHandler.GetDepartment)
					depts.PUT("/:id", rbacMiddleware.RequirePermission("departments.update"), departmentHandler.UpdateDepartment)
					depts.DELETE("/:id", rbacMiddleware.RequirePermission("departments.delete"), departmentHandler.DeleteDepartment)
				}

				// Services
				services := system.Group("/services")
				{
					services.POST("", rbacMiddleware.RequirePermission("services.create"), medicalServiceHandler.CreateService)
					services.GET("", rbacMiddleware.RequirePermission("services.view"), medicalServiceHandler.ListServices)
					services.GET("/:id", rbacMiddleware.RequirePermission("services.view"), medicalServiceHandler.GetService)
					services.PUT("/:id", rbacMiddleware.RequirePermission("services.update"), medicalServiceHandler.UpdateService)
				}

				// Audit Logs
				audit := system.Group("/audit-logs")
				{
					audit.GET("", rbacMiddleware.RequirePermission("audit.view"), auditLogHandler.ListLogs)
				}
			}
		}
	}
}
